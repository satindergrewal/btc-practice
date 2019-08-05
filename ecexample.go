package main

import (
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"crypto/sha256"
	"fmt"
)

func main() {
	privateKey, err := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
	fmt.Printf("%s\n\n", privateKey)
	fmt.Printf("%s\n\n", &privateKey.PublicKey)
	if err != nil {
		panic(err)
	}

	msg := "hello, world"
	hash := sha256.Sum256([]byte(msg))

	r, s, err := ecdsa.Sign(rand.Reader, privateKey, hash[:])
	if err != nil {
		panic(err)
	}
	fmt.Printf("signature: (0x%x, 0x%x)\n", r, s)

	valid := ecdsa.Verify(&privateKey.PublicKey, hash[:], r, s)
	fmt.Println("signature verified:", valid)
}
