/*
 * The following code is a demo.  Edge cases and error checking are
 * intentionally omitted where they might otherwise distract us from the
 * core ideas.
 *     "Before we can understand how something can go wrong, we must
 *      learn how it can go right."
 */
package main

import (
	"crypto/sha256"
	"fmt"
	"math/big"

	"golang.org/x/crypto/ripemd160"
)

type Point struct {
	x, y *big.Int
}

func NewPoint() Point {
	return Point{new(big.Int), new(big.Int)}
}

func (P Point) String() string {
	return fmt.Sprintf("Point(%d, %d)", P.x, P.y)
}

func (P Point) Cmp(Q Point) bool {
	return P.x.Cmp(Q.x) == 0 && P.y.Cmp(Q.y) == 0
}

func (P Point) Set(Q Point) Point {
	P.x.Set(Q.x)
	P.y.Set(Q.y)
	return P
}

var p = new(big.Int)
var identity = NewPoint()

/*
 * Elliptic curve point addition.  Unneeded side-cases are omitted for
 * simplicity.  See:
 *     https://en.wikipedia.org/wiki/Elliptic_curve_point_multiplication#Point_addition
 *     https://stackoverflow.com/a/31089415
 *     https://en.wikipedia.org/wiki/Modular_multiplicative_inverse#Using_Euler's_theorem
 *     https://crypto.stanford.edu/pbc/notes/elliptic/explicit.html
 */
func (R Point) ECPointAdd(P, Q Point) Point {
	if P.Cmp(identity) {
		return R.Set(Q)
	}
	s := new(big.Int)
	if P.Cmp(Q) {
		s.Mul(big.NewInt(2), P.y)
		s.ModInverse(s, p)
		s.Mul(s, big.NewInt(3))
		s.Mul(s, P.x)
		s.Mul(s, P.x)
	} else {
		s.Sub(Q.x, P.x)
		s.Mod(s, p)
		s.ModInverse(s, p)
		s.Mul(s, new(big.Int).Sub(Q.y, P.y))
	}
	x := new(big.Int)
	y := new(big.Int)
	x.Mul(s, s)
	x.Sub(x, P.x)
	x.Sub(x, Q.x)
	x.Mod(x, p)
	y.Sub(P.x, x)
	y.Mul(y, s)
	y.Sub(y, P.y)
	y.Mod(y, p)
	return R.Set(Point{x, y})
}

/*
 * Elliptic curve point multiplication.  This is an implimentation of
 * the Double-and-add algorithm with increasing index described here:
 *     https://en.wikipedia.org/wiki/Elliptic_curve_point_multiplication#Double-and-add
 */
func (Q Point) ECPointMul(d *big.Int, P Point) Point {
	N := NewPoint()
	N.Set(P)
	Q.Set(identity)
	for i := 0; i < d.BitLen(); i++ {
		if d.Bit(i) == 1 {
			Q.ECPointAdd(Q, N)
		}
		N.ECPointAdd(N, N)
	}
	return Q
}

/*
 * The compressed serialization of the public key.  See:
 *     Mastering Bitcoin, pages 73-75.
 *     https://www.ntirawen.com/2019/03/bitcoin-compressed-and-uncompressed.html
 *     https://www.secg.org/sec2-v2.pdf - Section 2.4.1.
 *     https://www.secg.org/sec1-v2.pdf - Section 2.3.3.
 */
func (R Point) Serialize() []byte {
	var s [33]byte
	s[0] = byte(2 + R.y.Bit(0))
	x := R.x.Bytes()
	copy(s[33-len(x):], x)
	return s[:]
}

/*
 * RIPEMD-160 hash.
 */
func r160(data []byte) []byte {
	d := ripemd160.New()
	d.Write(data)
	return d.Sum(nil)
}

/*
 * SHA-256 hash.
 */
func s256(data []byte) []byte {
	d := sha256.New()
	d.Write(data)
	return d.Sum(nil)
}

/*
 * Encode the data in Bitcoin's Base58 format.  See:
 *     https://en.bitcoin.it/wiki/Base58Check_encoding#Base58_symbol_chart
 */
func base58(data []byte) string {
	zero := new(big.Int)
	base := big.NewInt(58)
	symbol := "123456789ABCDEFGHJKLMNPQRSTUVWXYZabcdefghijkmnopqrstuvwxyz"

	var result []byte
	x := new(big.Int).SetBytes(data)
	r := new(big.Int)
	for x.Cmp(zero) != 0 {
		x.DivMod(x, base, r)
		result = append(result, symbol[r.Int64()])
	}
	for i := 0; data[i] == 0; i++ {
		result = append(result, symbol[0])
	}
	for i, j := 0, len(result)-1; i < j; i, j = i+1, j-1 {
		result[i], result[j] = result[j], result[i]
	}
	return string(result)
}

func main() {
	fmt.Println("Example from Mastering Bitcoin, pages 69-70.")

	/*
	 * Secp256k1 parameters.  See:
	 *     https://en.bitcoin.it/wiki/Secp256k1
	 *     https://www.secg.org/sec2-v2.pdf - Section 2.4.1.
	 */
	p.SetBit(p, 256, 1).Sub(p, big.NewInt(1<<32+977))
	fmt.Println("\tp:")
	fmt.Println(p)

	G := NewPoint()
	G.x.SetString("79be667ef9dcbbac55a06295ce870b07029bfcdb2dce28d9"+
		"59f2815b16f81798", 16)
	G.y.SetString("483ada7726a3c4655da4fbfc0e1108a8fd17b448a6855419"+
		"9c47d08ffb10d4b8", 16)
	fmt.Println("\tG:")
	fmt.Println(G)

	/*
	 * A private key must be a whole number from 1 to 0xffffffff...
	 * fffffffffffffffffffffffebaaedce6af48a03bbfd25e8cd0364140, one
	 * less than the order of the base point (or "generator point")
	 * G.  See:
	 *     https://en.bitcoin.it/wiki/Private_key//Range_of_valid_ECDSA_private_keys
	 */
	var privateKey = new(big.Int)
	privateKey.SetString("038109007313a5807b2eccc082c8c3fbb988a973"+
		"cacf1a7df9ce725c31b14776", 16)
	fmt.Println("\tprivateKey:")
	fmt.Println(privateKey)

	/*
	 * Generate the public key.  See:
	 *     Mastering Bitcoin, page 63.
	 */
	var publicKey = NewPoint()
	publicKey.ECPointMul(privateKey, G)
	fmt.Println("\tpublicKey:")
	fmt.Println(publicKey)

	serializedPublicKey := publicKey.Serialize()
	fmt.Println("\tserializedPublicKey:")
	fmt.Printf("%x\n", serializedPublicKey)

	publicKeyHash := r160(s256(serializedPublicKey))
	fmt.Println("\tpublicKeyHash:")
	fmt.Printf("%x\n", publicKeyHash)

	/*
	 * Calculate the checksum needed for Bitcoin's Base58Check
	 * format.  See:
	 *     Mastering Bitcoin, page 58
	 *     https://en.bitcoin.it/wiki/Technical_background_of_version_1_Bitcoin_addresses#How_to_create_Bitcoin_Address - Steps 5-7.
	 */
	version := []byte{0}
	versionPlusHash := append(version, publicKeyHash...)
	checksum := s256(s256(versionPlusHash))[:4]
	fmt.Println("\tchecksum:")
	fmt.Printf("%x\n", checksum)

	/*
	 * A Bitcoin address is just the public key hash encoded in
	 * Bitcoin's Base58Check format.  See:
	 *     Mastering Bitcoin, page 66.
	 */
	address := base58(append(versionPlusHash, checksum...))
	fmt.Println("\taddress:")
	fmt.Println(address)
}
