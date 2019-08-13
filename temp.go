package main

import (
	"fmt"
	"log"
	"math/big"
)

func primesLessThan(n *big.Int) (primes []big.Int) {
	var one big.Int
	one.SetInt64(1)
	var i big.Int
	i.SetInt64(1)
	for i.Cmp(n) < 0 {
		var result big.Int
		result.Set(&i)
		fmt.Println(result.String())
		primes = append(primes, result)
		i.Add(&i, &one)
	}
	return
}

func main() {
	//primes := primesLessThan(big.NewInt(5))
	//for _, p := range primes {
	//fmt.Println(p.String())
	//}

	//var P []big.Int
	//P[0].SetInt64(1)
	x := new(big.Int)
	x, ok := x.SetString("55066263022277343669578718895168534326250603453777594175500187360389116729240", 10)
	if !ok {
		log.Fatalf("big Int value did not set")
		return
	}
	fmt.Println(x)

	y := new(big.Int)
	y, ok = y.SetString("32670510020758816978083085130507043184471273380659243275938904335757337482424", 10)
	if !ok {
		log.Fatalf("big Int value did not set")
		return
	}
	fmt.Println(y)

	var P [2]big.Int
	//P[0] = *P[0].SetInt64(1)
	P[0] = *x
	P[1] = *y
	fmt.Printf("%t\n", P)

	if P[0].Cmp(x) != 0 {
		panic("p is not equals to p_hex")
		//return errors.New("The point does not lie on the elliptic curve")
	}
}
