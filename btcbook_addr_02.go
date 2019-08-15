package main

// BTC Address generation pracitce code
// converting python code example from Gareth's gareth_file04.py

import (
	"fmt"
	"log"
	"math/big"
)

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

	// Secp256k1 parameters.  See:
	//     https://en.bitcoin.it/wiki/Secp256k1
	//     https://www.secg.org/sec2-v2.pdf - Section 2.4.1.
	p := big.NewInt(0)
	p.Exp(big.NewInt(2), big.NewInt(256), nil)
	p.Sub(p, big.NewInt(1<<32+1<<9+1<<8+1<<7+1<<6+1<<4+1))
	//fmt.Println(p)

	/*G = Point(
	0x79be667ef9dcbbac55a06295ce870b07029bfcdb2dce28d959f2815b16f81798,
	0x483ada7726a3c4655da4fbfc0e1108a8fd17b448a68554199c47d08ffb10d4b8,
	)
	*/

}
