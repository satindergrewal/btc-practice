package main

// BTC Address generation pracitce code
// converting python code example from Gareth's gareth_file04.py

import (
	"fmt"
	"log"
	"math/big"
)

type Point struct {
	x big.Int
	y big.Int
}

func ec_p() big.Int {
	// Secp256k1 parameters.  See:
	//     https://en.bitcoin.it/wiki/Secp256k1
	//     https://www.secg.org/sec2-v2.pdf - Section 2.4.1.
	// Defining p
	p := big.NewInt(0)
	p.Exp(big.NewInt(2), big.NewInt(256), nil)
	p.Sub(p, big.NewInt(1<<32+1<<9+1<<8+1<<7+1<<6+1<<4+1))
	//fmt.Println(p)
	return *p
}

func ec_G() Point {
	// Defining G
	x := new(big.Int)
	x, ok := x.SetString("79be667ef9dcbbac55a06295ce870b07029bfcdb2dce28d959f2815b16f81798", 16)
	if !ok {
		log.Fatalf("big Int value did not set")
	}
	y := new(big.Int)
	y, ok = y.SetString("483ada7726a3c4655da4fbfc0e1108a8fd17b448a68554199c47d08ffb10d4b8", 16)
	if !ok {
		log.Fatalf("big Int value did not set")
	}
	//fmt.Println(x)
	//fmt.Println(y)

	var G Point
	G.x = *x
	G.y = *y
	//fmt.Printf("%d\n", G.x)
	return G
}

// Elliptic curve point addition.  Unneeded side-cases are omitted for
// simplicity.  See:
//     https://en.wikipedia.org/wiki/Elliptic_curve_point_multiplication#Point_addition
//     https://stackoverflow.com/a/31089415
//     https://en.wikipedia.org/wiki/Modular_multiplicative_inverse#Using_Euler's_theorem
//     https://crypto.stanford.edu/pbc/notes/elliptic/explicit.html
func ec_point_add(P, Q *Point) Point {
	//fmt.Println(P)
	//fmt.Println(Q)

	if big.NewInt(0).Cmp(&P.x) == 0 && big.NewInt(0).Cmp(&P.y) == 0 {
		return *Q
	}

	p := *big.NewInt(0)
	p = ec_p()
	//fmt.Println("value of p is: ", p)

	slope := big.NewInt(0)

	if P == Q {
		// slope = 3*P.x**2 * pow(2*P.y, p-2, p)  # 3Px^2 / 2Py
		pxpow2 := big.NewInt(0).Exp(&P.x, big.NewInt(2), nil)   // Px^2
		threepxpow2 := big.NewInt(0).Mul(big.NewInt(3), pxpow2) //3Px^2
		pymul2 := big.NewInt(0).Mul(big.NewInt(2), &P.y)        // 2*P.y
		psub2 := big.NewInt(0).Sub(&p, big.NewInt(2))           // p-2
		expo := big.NewInt(0).Exp(pymul2, psub2, &p)            // pow(2*P.y, p-2, p)
		slope = big.NewInt(0).Mul(threepxpow2, expo)            // 3Px^2 / 2Py
		//fmt.Println("\n\nSLOPE 1: ", slope)
	} else {
		// slope = (Q.y - P.y) * pow(Q.x - P.x, p-2, p)  # (Qy - Py) / (Qx - Px)
		qysubpy := big.NewInt(0).Sub(&Q.y, &P.y)      // Qy - Py
		qxsubpx := big.NewInt(0).Sub(&Q.x, &P.x)      // Qx - Px
		psub2 := big.NewInt(0).Sub(&p, big.NewInt(2)) // p-2
		expo := big.NewInt(0).Exp(qxsubpx, psub2, &p) // pow(Q.x - P.x, p-2, p)
		slope = big.NewInt(0).Mul(qysubpy, expo)      // (Q.y - P.y) * pow(Q.x - P.x, p-2, p)
		//fmt.Println("\n\nSLOPE 2: ", slope)
	}

	var R Point
	rx := big.NewInt(0).Exp(slope, big.NewInt(2), nil) // slope^2
	rx = big.NewInt(0).Sub(rx, &P.x)                   // slope^2 - Px
	rx = big.NewInt(0).Sub(rx, &Q.x)                   // slope^2 - Px - Qx
	R.x = *big.NewInt(0).Exp(rx, &p, &p)               // Mod value

	ry := big.NewInt(0).Sub(&P.x, &R.x) // (Px - Rx)
	ry = big.NewInt(0).Mul(slope, ry)   // slope*(Px - Rx)
	ry = big.NewInt(0).Sub(ry, &P.y)    // slope*(Px - Rx) - Py
	R.y = *big.NewInt(0).Mod(ry, &p)    // Mod value
	//fmt.Println(&R)

	return R
}

// Elliptic curve point multiplication.  This is an implimentation of the
// Double-and-add algorithm with increasing index described here:
//     https://en.wikipedia.org/wiki/Elliptic_curve_point_multiplication#Double-and-add
func ec_point_multiply(d *big.Int, P *Point) Point {
	N := &P
	var Q Point
	Q.x = *big.NewInt(0)
	Q.y = *big.NewInt(0)

	fmt.Println("Value of d: ", d.Bit(0))
	fmt.Println("BitLen of d: ", d.BitLen())

	//fmt.Println(d & big.NewInt(1))

	for i := 0; i <= d.BitLen(); i++ {
		//fmt.Println(i, d.Bit(i))
		if d.Bit(i) == 1 {
			//fmt.Println(i, d.Bit(i))
			Q = ec_point_add(&Q, *N)
			//fmt.Println(&Q)
		} else {
			//*N = ec_point_add(*N, *N)
			fmt.Println("N is: ", *N)
		}
	}

	/*while d:
	  if d & 1:
	      Q = ec_point_add(Q, N)
	  N = ec_point_add(N, N)
	  d >>= 1*/
	return Q
}

func main() {
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
	private_key, ok := private_key.SetString("038109007313a5807b2eccc082c8c3fbb988a973cacf1a7df9ce725c31b14776", 16)
	if !ok {
		log.Fatalf("big Int value did not set")
		//return errors.New("big Int value did not set")
	}
	fmt.Printf("private_key: %d\n", private_key)

	var G Point
	G = ec_G()
	fmt.Printf("Gx: %d\n", &G.x)
	fmt.Printf("Gy: %d\n", &G.y)

	var G2 Point
	G2 = ec_point_add(&G, &G)
	//fmt.Println("Value of G2: ", G2)
	fmt.Printf("G2x: %d\n", &G2.x)
	fmt.Printf("G2y: %d\n", &G2.y)

	//theeg := ec_point_add(&G, &P)
	//fmt.Printf("theeg.x: %d\n", &theeg.x)
	//fmt.Printf("theeg.y: %d\n", &theeg.y)

	var Pmul Point
	Pmul = ec_point_multiply(private_key, &G2)
	fmt.Printf("\n\n%t\n", &Pmul.x)
	if big.NewInt(0).Cmp(&Pmul.x) == 0 && big.NewInt(0).Cmp(&Pmul.y) == 0 {
		fmt.Printf("Pmul.x is not zero: %d\n", &Pmul.x)
	} else {
		fmt.Printf("Pmul.x is zero: %d\n", &Pmul.x)
	}
}
