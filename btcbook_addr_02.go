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

}
