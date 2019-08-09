package main

import (
	//"crypto/ecdsa"
	//"crypto/elliptic"
	//"crypto/rand"
	//"crypto/sha256"
	"fmt"
	"math/big"
	//"math"
	"errors"
	"log"
)

func ec_valid(x, y *big.Int) error {
	// The characteristic of secp256k1; the order of the corresponding finite field.
	// defined big int, as the code was throwing error "constant 115...663 overflows int"
	p := new(big.Int)
	p, ok := p.SetString("115792089237316195423570985008687907853269984665640564039457584007908834671663", 10)
	if !ok {
		return errors.New("big Int value did not set")
	}
	//fmt.Println(p)

	// convert big int value to hex and store in a variable and print it on screen.
	//phx := fmt.Sprintf("%X", p)
	//fmt.Println(phx)

	// Check that the value of `p` given in the book matches the hexadecimal value
	// given on en.bitcoin.it.
	// source of hex number: https://en.bitcoin.it/wiki/Secp256k1
	// converting the hex value of big int number and converitng it back to big int number.
	p_hex := new(big.Int)
	p_hex, _ = p_hex.SetString("fffffffffffffffffffffffffffffffffffffffffffffffffffffffefffffc2f", 16)
	//fmt.Println(p_hex)

	if p.Cmp(p_hex) != 0 {
		//panic("p is not equals to p_hex")
		return errors.New("The point does not lie on the elliptic curve")
	}

	// Check that the value of `p` given in the book matches the mathematical
	// expression given in the NIST spec.
	// Use this example by Gareth: https://play.golang.org/p/gazNAmBasle
	// Or use this second better, smaller version of the same example: https://play.golang.org/p/WtK1OqNA8LJ
	p_math := big.NewInt(0)
	p_math.Exp(big.NewInt(2), big.NewInt(256), nil)
	p_math.Sub(p_math, big.NewInt(1<<32 + 1<<9 + 1<<8 + 1<<7 + 1<<6 + 1<<4 + 1))
	//fmt.Println(p_math)


	// Example point on the curve 'secp256k1'.
	/*x := new(big.Int)
	x, ok = x.SetString("55066263022277343669578718895168534326250603453777594175500187360389116729240", 10)
	if !ok {
		fmt.Printf("NOT WORKING!")
		return
	}
	y := new(big.Int)
	y, ok = y.SetString("32670510020758816978083085130507043184471273380659243275938904335757337482424", 10)
	if !ok {
		fmt.Printf("NOT WORKING!")
		return
	}*/
	//fmt.Println(x)
	//fmt.Println(y)
    
    // Check that the above point actually lies on the elliptic curve
	//     y^2 = x^3 + ax + b.
	//		x = 0
	// 		b = 7
	// can just use p instead of nil to store get mod value as output
	left_side := big.NewInt(0).Exp(y, big.NewInt(2), p)
	// left_side.Mod(left_side, p)
	//fmt.Println("y^2 = ",left_side)

	//right_side = (x**3 + 7) % p
	right_side := big.NewInt(0).Exp(x, big.NewInt(3), nil)
	right_side.Add(right_side, big.NewInt(7))
	right_side.Mod(right_side, p)
	//fmt.Println("x^3 + ax + b = ",right_side)

	if left_side.Cmp(right_side) != 0 {
		//panic("left_side is not equals to right_side")
		return errors.New("The point does not lie on the elliptic curve")
	} else {
		return nil
	}
}

func main() {
	x := new(big.Int)
	x, ok := x.SetString("55066263022277343669578718895168534326250603453777594175500187360389116729240", 10)
	if !ok {
		log.Fatalf("big Int value did not set")
		return
	}
	y := new(big.Int)
	y, ok = y.SetString("32670510020758816978083085130507043184471273380659243275938904335757337482424", 10)
	if !ok {
		log.Fatalf("big Int value did not set")
		return
	}

	err := ec_valid(x, y)
	if err != nil {
		fmt.Println(err)
	} else {
		fmt.Println("Yes the code works!")
	}
}
