package main

// BTC Address generation pracitce code
// converting python code example from Gareth's gareth_file04.py

import (
	"bytes"
	"crypto/ed25519"
	"crypto/sha256"
	"encoding/base32"
	"fmt"
	"strings"

	// "log"
	"math/big"

	"golang.org/x/crypto/ripemd160"
	"golang.org/x/crypto/sha3"
)

type Point struct {
	x, y *big.Int
}

// Secp256k1 parameters.  See:
//     https://en.bitcoin.it/wiki/Secp256k1
//     https://www.secg.org/sec2-v2.pdf - Section 2.4.1.
// Defining p
func ec_p() *big.Int {
	p := big.NewInt(0)
	p.Exp(big.NewInt(2), big.NewInt(256), nil)
	p.Sub(p, big.NewInt(1<<32+1<<9+1<<8+1<<7+1<<6+1<<4+1))
	//fmt.Println(p)
	return p
}

// Defining G
func ec_G() Point {
	var G = Point{new(big.Int), new(big.Int)}
	G.x.SetString("79be667ef9dcbbac55a06295ce870b07029bfcdb2dce28d959f2815b16f81798", 16)
	G.y.SetString("483ada7726a3c4655da4fbfc0e1108a8fd17b448a68554199c47d08ffb10d4b8", 16)
	//fmt.Println(G.x)
	//fmt.Println(G.y)
	return G
}

// Elliptic curve point addition.  Unneeded side-cases are omitted for
// simplicity.  See:
//     https://en.wikipedia.org/wiki/Elliptic_curve_point_multiplication#Point_addition
//     https://stackoverflow.com/a/31089415
//     https://en.wikipedia.org/wiki/Modular_multiplicative_inverse#Using_Euler's_theorem
//     https://crypto.stanford.edu/pbc/notes/elliptic/explicit.html
func (R *Point) ec_point_add(P, Q *Point) *Point {
	//fmt.Println(P)
	//fmt.Println("Qx: ", Q.x)
	//fmt.Println("Qy: ", Q.y)

	if big.NewInt(0).Cmp(P.x) == 0 && big.NewInt(0).Cmp(P.y) == 0 {
		R.x.Set(Q.x)
		R.y.Set(Q.y)
		return R
	}

	p := big.NewInt(0)
	p = ec_p()
	//fmt.Println("value of p is: ", p)

	slope := big.NewInt(0)

	if P.x.Cmp(Q.x) == 0 && P.y.Cmp(Q.y) == 0 {
		// slope = 3*P.x**2 * pow(2*P.y, p-2, p)  # 3Px^2 / 2Py
		pxpow2 := big.NewInt(0).Exp(P.x, big.NewInt(2), nil)    // Px^2
		threepxpow2 := big.NewInt(0).Mul(big.NewInt(3), pxpow2) //3Px^2
		pymul2 := big.NewInt(0).Mul(big.NewInt(2), P.y)         // 2*P.y
		psub2 := big.NewInt(0).Sub(p, big.NewInt(2))            // p-2
		expo := big.NewInt(0).Exp(pymul2, psub2, p)             // pow(2*P.y, p-2, p)
		slope = big.NewInt(0).Mul(threepxpow2, expo)            // 3Px^2 / 2Py
		//fmt.Println("\n\nSLOPE 1: ", slope)
	} else {
		// slope = (Q.y - P.y) * pow(Q.x - P.x, p-2, p)  # (Qy - Py) / (Qx - Px)
		//fmt.Printf("Qy %d\nPy %d\n", Q.y, P.y)
		//fmt.Printf("Qx %d\nPx %d\n", Q.x, P.x)
		qysubpy := big.NewInt(0).Sub(Q.y, P.y) // Qy - Py
		//fmt.Println("qysubpy", qysubpy)
		qxsubpx := big.NewInt(0).Sub(Q.x, P.x) // Qx - Px
		//fmt.Println("qxsubpx", qxsubpx)
		//qxsubpxmodp := big.NewInt(0).Mod(qxsubpx, p)     // (Qx - Px)%p
		psub2 := big.NewInt(0).Sub(p, big.NewInt(2)) // p-2
		//fmt.Println("psub2", psub2)
		expo := big.NewInt(0).Exp(qxsubpx, psub2, p) // pow(Q.x - P.x, p-2, p)
		//fmt.Println("expo", expo)
		slope = big.NewInt(0).Mul(qysubpy, expo) // (Q.y - P.y) * pow(Q.x - P.x, p-2, p)
		//fmt.Println("\n\nSLOPE 2: ", slope)
	}

	rx := big.NewInt(0).Exp(slope, big.NewInt(2), nil) // slope^2
	rx = big.NewInt(0).Sub(rx, P.x)                    // slope^2 - Px
	rx = big.NewInt(0).Sub(rx, Q.x)                    // slope^2 - Px - Qx
	rx.Exp(rx, p, p)                                   // Mod value

	//fmt.Println("P.x", P.x)
	//fmt.Println("rx", rx)
	ry := big.NewInt(0).Sub(P.x, rx) // (Px - Rx)
	//fmt.Println("ry", ry)
	ry = big.NewInt(0).Mul(slope, ry) // slope*(Px - Rx)
	ry = big.NewInt(0).Sub(ry, P.y)   // slope*(Px - Rx) - Py
	ry.Mod(ry, p)                     // Mod value
	//fmt.Println(&R)

	R.x.Set(rx)
	R.y.Set(ry)
	return R
}

// Elliptic curve point multiplication.  This is an implimentation of the
// Double-and-add algorithm with increasing index described here:
//     https://en.wikipedia.org/wiki/Elliptic_curve_point_multiplication#Double-and-add
func (Q *Point) ec_point_multiply(d *big.Int, P *Point) Point {
	N := &Point{new(big.Int), new(big.Int)}
	N.x.Set(P.x)
	N.y.Set(P.y)

	//var Q Point
	Q.x = big.NewInt(0)
	Q.y = big.NewInt(0)

	//fmt.Println("Value of d: ", d)
	//fmt.Println("Value of d.Bit(0): ", d.Bit(0))
	//fmt.Println("BitLen of d: ", d.BitLen())

	for i := 0; i <= d.BitLen(); i++ {
		//fmt.Println(i, d.Bit(i))
		if d.Bit(i) == 1 {
			Q.ec_point_add(Q, N)
			//fmt.Println(i, d.Bit(i), Q.x, Q.y)
		}
		N.ec_point_add(N, N)
		//fmt.Println("N is: ", N.x, N.y)
	}

	//fmt.Println("ec multiply return: ", Q.x, Q.y)
	return *Q
}

// Following the exmple of Figure 4-7 of Mastering Bitcoin
func (R Point) Serialize() []byte {
	b := R.x.Bytes()
	fmt.Println("Check if Y is Even or Odd: ")
	fmt.Println("R.y.Mod(big.NewInt(2))", R.y.Mod(R.y, big.NewInt(2)))

	// If the length of Public Key x bytes is lesser than 32 bytes, we need to add
	// the required remaining bytes to the Public Key x.
	// Check the length of byte slice for R.x
	if length := 0; len(b) < 32 { // temp variable to store length value for public key x
		fmt.Println("Public Key x Byte length", len(b))
		length = 32 - len(b)
		fmt.Printf("Byte length of public key x is short by %d\n", length)
		addbytes := make([]byte, length)
		fmt.Println("addbytes", addbytes)
		b = append(addbytes, b...)
		fmt.Println("b byte length after addbytes", len(b))
	}

	/*
			// Even odd number checking
			if(n%2==0){
		        fmt.Println(n,"is Even number")
		    }else{
		        fmt.Println(n,"is Odd number")
		    }
	*/
	if R.y.Mod(R.y, big.NewInt(2)).Cmp(big.NewInt(0)) == 0 {
		fmt.Println(R.y, "is Even number")
		fmt.Println("R.x", R.x)
		fmt.Println("R.x Bytes", b)
		fmt.Println("R.x Bytes length", len(b))
		fmt.Println("byte 02", []byte{02})
		a := append([]byte{02}, b...)
		fmt.Println("a", a)
		return a
	} else {
		fmt.Println(R.y, "is Odd number")
		a := append([]byte{03}, b...)
		fmt.Println("a", a)
		return a
	}
}

/*
 * RIPEMD-160 hash.
 */
func r160(data []byte) []byte {
	h := ripemd160.New()
	h.Write(data)
	return h.Sum(nil)
}

/*
 * SHA-256 hash.
 */
func s256(data []byte) []byte {
	h := sha256.New()
	h.Write(data)
	return h.Sum(nil)
}

/*
 * Encode the data in Bitcoin's Base58 format.  See:
 *     https://en.bitcoin.it/wiki/Base58Check_encoding#Base58_symbol_chart
 */
func base58(data []byte) string {
	base := big.NewInt(58)                                                        // set new big int for dividing data bytes with 58
	remainder := new(big.Int)                                                     // to store remainder of the data bytes
	const alphabet = "123456789ABCDEFGHJKLMNPQRSTUVWXYZabcdefghijkmnopqrstuvwxyz" // base 58 character set

	//fmt.Println("Length of base58 characters: ", len(alphabet))
	//fmt.Printf("Recieved data bytes: %x\n", data)

	x := new(big.Int).SetBytes(data) // convert bytes to big integer
	//fmt.Println("Bytes to big integer: ", x)
	var output_string []byte // To store base58 coverted value to a temp variable

	for x.Cmp(big.NewInt(0)) != 0 {
		x.DivMod(x, base, remainder)
		//fmt.Println("X", x, "\nRemainder", remainder)
		//fmt.Printf("alphabet[%d]: %s\n", remainder.Int64(), string(alphabet[remainder.Int64()]))
		output_string = append(output_string, alphabet[remainder.Int64()]) // Appending the value stored at position alphabet[remainder] to output_string
	}
	//fmt.Println("output_string value after first loop: ", string(output_string))

	for i := 0; data[i] == 0; i++ {
		output_string = append(output_string, alphabet[0])
	}
	//fmt.Println("output_string value after second loop: ", string(output_string))

	// Reverse the array.
	for i, j := 0, len(output_string)-1; i < j; i, j = i+1, j-1 {
		output_string[i], output_string[j] = output_string[j], output_string[i]
	}
	return string(output_string)
}

func main() {
	// btcVersionByte := []byte{0x0}
	// btcPrivKeyVersionByte := []byte{0x80}
	kmdVersionByte := []byte{0x3C}
	kmdPrivKeyVersionByte := []byte{0xBC}
	// Please note that the following code is a demo.  Edge cases and error
	// checking are intentionally omitted where they might otherwise distract
	// us from the core ideas.
	//     "Before we can understand how something can go wrong, we must
	//      learn how it can go right."

	fmt.Println("Example from Mastering Bitcoin, pages 69-70.")

	// A private key must be a whole number from 1 to
	//     0xfffffffffffffffffffffffffffffffebaaedce6af48a03bbfd25e8cd0364140,
	// one less than the order of the base point (or "generator point") G.
	// See:
	//     https://en.bitcoin.it/wiki/Private_key//Range_of_valid_ECDSA_private_keys
	private_key := new(big.Int)

	// Mastering Bitcoin example privkey, which has odd public key x value
	//private_key, ok := private_key.SetString("038109007313a5807b2eccc082c8c3fbb988a973cacf1a7df9ce725c31b1477a", 16)

	passStr := "myverysecretandstrongpassphrase_noneabletobrute" // Password string or Pass phrase
	hashByte := s256([]byte(passStr))

	// The following is the method used by iguana password hashing algorigthm, as shown in existing seed to address/WIF generate examples of PHP and JavaScript
	// See:
	//		https://github.com/pbca26/komodolib-js/blob/master/src/keys.js#L94
	//		https://github.com/DeckerSU/komodo_scripts/blob/master/genkomodo.php#L136
	// fmt.Println(hashByte)
	// fmt.Println(hashByte[0] & 248)
	// fmt.Println(hashByte[31] & 127)
	// fmt.Println(hashByte[31] | 64)
	hashByte[0] &= 248
	// fmt.Println(hashByte[0])
	hashByte[31] &= 127
	// fmt.Println(hashByte[31])
	hashByte[31] |= 64
	// fmt.Println(hashByte[31])
	passHash := hashByte

	/// The following passHash is used by the Bitcoin hasing algorithm, as explained in Mastering Bitcoin and othe resources
	// passHash := s256(s256(hashByte))      // convert passphrase to bytes/uint8 and double hash it with SHA256
	private_key = private_key.SetBytes(passHash) // Setting passrase hash with SetBytes
	fmt.Println("Password/Passphrase: ", passStr)
	fmt.Printf("password hash: %x\n", passHash)

	// Mastering Bitcoin example privkey, which has even public key x value
	// private_key, ok := private_key.SetString("038109007313a5807b2eccc082c8c3fbb988a973cacf1a7df9ce725c31b14776", 16)
	// if !ok {
	// 	log.Fatalf("big Int value did not set")
	// 	//return errors.New("big Int value did not set")
	// }
	fmt.Printf("private_key: %d\n", private_key)

	var G Point
	G = ec_G()
	fmt.Printf("Gx: %d\n", G.x)
	fmt.Printf("Gy: %d\n", G.y)

	var publicKey Point
	publicKey.ec_point_multiply(private_key, &G)
	fmt.Printf("\npublicKey.x %d\n", publicKey.x)
	fmt.Printf("publicKey.y %d\n", publicKey.y)

	serializedPublicKey := publicKey.Serialize()
	fmt.Printf("serializedPublicKey: %x\n", serializedPublicKey)

	publicKeyHash := r160(s256(serializedPublicKey))
	fmt.Printf("publicKeyHash: %x\n", publicKeyHash)
	publicKeyHashSHA256 := s256(publicKeyHash)

	fmt.Printf("publicKeyHashSHA256: %x\n", publicKeyHashSHA256)
	fmt.Println("publicKeyHashSHA256 Len:", len(publicKeyHashSHA256))

	onionPrivKey := ed25519.NewKeyFromSeed([]byte(publicKeyHashSHA256))
	fmt.Printf("\nonionPrivKey: %x\n", onionPrivKey)
	onionPubKey := onionPrivKey.Public()
	fmt.Printf("onionPubKey: %x\n", onionPubKey)

	// checksum = H(".onion checksum" || pubkey || version)
	var onionChecksumBytes bytes.Buffer
	onionChecksumBytes.Write([]byte(".onion checksum"))
	onionChecksumBytes.Write([]byte(onionPubKey.(ed25519.PublicKey)))
	onionChecksumBytes.Write([]byte{0x03})
	onionChecksum := sha3.Sum256(onionChecksumBytes.Bytes())
	fmt.Printf("onionChecksum %x\n", onionChecksum)

	// onion_address = base32(pubkey || checksum || version)
	var onionAddressBytes bytes.Buffer
	onionAddressBytes.Write([]byte(onionPubKey.(ed25519.PublicKey)))
	onionAddressBytes.Write([]byte(onionChecksum[:2]))
	onionAddressBytes.Write([]byte{0x03})
	onionAddress := base32.StdEncoding.EncodeToString(onionAddressBytes.Bytes())

	fmt.Println("onionPubKey:", base32.StdEncoding.EncodeToString([]byte(onionPubKey.(ed25519.PublicKey))))
	fmt.Println("onionAddress:", strings.ToLower(onionAddress)+".onion", "\n\n")

	/*
	 * Calculate the checksum needed for Bitcoin's Base58Check
	 * format.  See:
	 *     Mastering Bitcoin, page 58
	 *     https://en.bitcoin.it/wiki/Technical_background_of_version_1_Bitcoin_addresses#How_to_create_Bitcoin_Address - Steps 5-7.
	 */
	version := kmdVersionByte
	fmt.Printf("Bitcoin version byte: %d\n", version)
	versionPlusHash := append(version, publicKeyHash...)
	fmt.Printf("bitcoin version + pubkey hash: %d\n", versionPlusHash)
	checksum := s256(s256(versionPlusHash))[:4]
	fmt.Printf("checksum: %x\n", checksum)

	/*
	 * A Bitcoin address is just the public key hash encoded in
	 * Bitcoin's Base58Check format.  See:
	 *     Mastering Bitcoin, page 66.
	 */
	address := base58(append(versionPlusHash, checksum...))
	fmt.Println("address:", address, "\n")

	/*
	 *	https://www.mobilefish.com/services/cryptocurrency/cryptocurrency.html#refPrivateKeyHex
	 */
	privKeyVersion := kmdPrivKeyVersionByte                              // version byte to add as prefix for private key
	privKeyPlusVersion := append(privKeyVersion, private_key.Bytes()...) // privkey version + privkey hash
	privKeyChecksum := s256(s256(privKeyPlusVersion))[:4]                // first 4 bytes of double hashed (privkey version + privkey hash)
	fmt.Printf("Bitcoin Private Key version byte: %d\n", privKeyVersion)
	fmt.Printf("bitcoin version + private key hash: %x   byte length: %d\n", privKeyPlusVersion, len(privKeyPlusVersion))
	fmt.Printf("privKeyChecksum: %x\n", privKeyChecksum)

	privKeyAddChecksumUncomp := append(privKeyPlusVersion, privKeyChecksum...) // privkey version + privkey hash + extra byte + last
	fmt.Printf("Privkey version + privkey hash + Uncompressed Checksum: %x   byte length: %d\n", privKeyAddChecksumUncomp, len(privKeyAddChecksumUncomp))

	wifUncompressed := base58(privKeyAddChecksumUncomp)
	fmt.Println("WIF (Uncompressed):", wifUncompressed, "\n")

	/*
	 * WIF Key (Compressed) need extra byte added to the end
	 * along with version + privkey hash + 01
	 */
	byteone := []byte{0x01}                                  // Extra byte 01 to add to PrivKey at the end
	privKeyVerByte := append(privKeyPlusVersion, byteone...) // privkey version + privkey hash + extra byte
	privKeyChecksumComp := s256(s256(privKeyVerByte))[:4]    // first 4 bytes of double hashed (privkey version + privkey hash + extra byte)
	fmt.Printf("byteone: %d\n", byteone)
	fmt.Printf("bitcoin version + private key hash + 01 byte: %x   byte length: %d\n", privKeyVerByte, len(privKeyVerByte))
	fmt.Printf("privKey Checksum for Compressed WIF: %x\n", privKeyChecksumComp)

	privKeyAddChecksum := append(privKeyVerByte, privKeyChecksumComp...) // privkey version + privkey hash + extra byte + last
	fmt.Printf("Privkey version + privkey hash + extra byte + Compressed Checksum: %x   byte length: %d\n", privKeyAddChecksum, len(privKeyAddChecksum))

	wifCompressed := base58(privKeyAddChecksum)
	fmt.Println("WIF (Compressed):", wifCompressed, "\n")
}
