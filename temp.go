package main

import (
	"fmt"
	"math/big"
)

func main() {
	a2 := big.NewInt(2)
	b256 := big.NewInt(256)
	b32 := big.NewInt(32)
	b9 := big.NewInt(9)
	//b8 := big.NewInt(8)
	//b7 := big.NewInt(7)
	//b6 := big.NewInt(6)
	//b4 := big.NewInt(4)
	//b0 := big.NewInt(0)
	sub1 := big.NewInt(0).Exp(a2, b256, nil)
	sub2 := big.NewInt(0).Exp(a2, b32, nil)
	sub3 := big.NewInt(0).Exp(a2, b9, nil)
	x := big.NewInt(0).Sub(sub1, sub2)
	y := big.NewInt(0).Sub(x, )
	fmt.Println(x)
}
