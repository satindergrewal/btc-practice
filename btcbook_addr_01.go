package main

import (
	//"crypto/ecdsa"
	//"crypto/elliptic"
	//"crypto/rand"
	//"crypto/sha256"
	"fmt"
	"math/big"
	"math"
	"strconv"
)

func FloatToString(input_num float64) string {
	// to convert a float number to a string
	return strconv.FormatFloat(input_num, 'f', 0, 64)
}

func RecursivePower(base int, exponent int) int {
	if exponent != 0 {
		return (base * RecursivePower(base, exponent-1))
	} else {
		return 1
	}
}

func FloatToBigInt(val float64) *big.Int {
	bigval := new(big.Float)
	bigval.SetFloat64(val)
	
	coin := new(big.Float)
	coin.SetInt(big.NewInt(1000000000000000000))
	
	bigval.Mul(bigval, coin)
	
	result := new(big.Int)
	bigval.Int(result)

	return result
}

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
	p_pow := math.Pow(2, 256) - math.Pow(2, 32) - math.Pow(2, 9) - math.Pow(2, 8) - math.Pow(2, 7) - math.Pow(2, 6) - math.Pow(2, 4) - math.Pow(2, 0)
	fmt.Printf("%.0f\n", p_pow)
	//p_math := new(big.Int)
	//p_math, _ = p_math.SetString(FloatToString(), 10)
	//p_math, _ = p_math.SetString(FloatToString(math.Pow(2, 256) - math.Pow(2, 32)), 10)
	//fmt.Printf("%.0f\n", p_math)
	//fmt.Println(p_math)
	//fmt.Println(p.Cmp(p_math))
	//assert p == p_math

	//fmt.Println(FloatToString(math.Pow(2, 256)))

	fmt.Printf("%f\n",math.Pow(2, 256))
	//pow := new(big.Int)
	//pow.SetInt(big.NewInt(4294967296))
	//pow = 4294967296
	pow := math.Pow(2, 256)-FloatToBigInt(float64(4294967296))
	fmt.Printf("%d\n", pow)
	fmt.Printf("%f\n",math.Pow(2, 9))
}
