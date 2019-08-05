package main

import (
	//"crypto/ecdsa"
	//"crypto/elliptic"
	//"crypto/rand"
	//"crypto/sha256"
	"fmt"
	"math/big"
)

func main() {
	// The characteristic of secp256k1; the order of the corresponding finite field.
	// defined big int, as the code was throwing error "constant 115...663 overflows int"
	p := new(big.Int)
	p, ok := p.SetString("115792089237316195423570985008687907853269984665640564039457584007908834671663", 10)
	if !ok {
		fmt.Printf("NOT WORKING!")
		return
	}
	fmt.Println(p)

	// convert big int value to hex and store in a variable and print it on screen.
	phx := fmt.Sprintf("%X", p)
	fmt.Println(phx)

	// Check that the value of `p` given in the book matches the hexadecimal value
	// given on en.bitcoin.it.
	// source of hex number: https://en.bitcoin.it/wiki/Secp256k1
	// converting the hex value of big int number and converitng it back to big int number.
	p_hex := new(big.Int)
	p_hex, _ = p_hex.SetString("fffffffffffffffffffffffffffffffffffffffffffffffffffffffefffffc2f", 16)
	fmt.Println(p_hex)

	if p.Cmp(p_hex) != 0 {
		panic("p is not equals to p_hex")
	}

	// Check that the value of `p` given in the book matches the mathematical
	// expression given in the NIST spec.
	//2 ^ 256
	p_math = 2**256 - 2**32 - 2**9 - 2**8 - 2**7 - 2**6 - 2**4 - 2**0
	//assert p == p_math

}
