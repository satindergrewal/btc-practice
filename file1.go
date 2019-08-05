package main

import (
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	//"crypto/sha256"
	"fmt"
	//"math/big"
)

func main() {
	//fmt.Printf("Hello World\n")

	privateKey, err := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
	fmt.Printf("%s\n\n", privateKey)
	//fmt.Printf("%s\n\n", &privateKey.PublicKey)
	if err != nil {
		panic(err)
	}


	fmt.Printf("%s\n\n", privateKey.D)

	dbyte := privateKey.D.Bytes()
	fmt.Printf("%v\n", dbyte)

}