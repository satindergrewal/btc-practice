package main

// Please note that the following code is a demo.  Edge cases and error
// checking are intentionally omitted where they might otherwise
// distract us from the core ideas.
//     "Before we can understand how something can go wrong, we must
//      learn how it can go right."

import (
    "fmt"
    "math/big"
)

type Point struct {
    x, y *big.Int
}

// Elliptic curve point addition.  Unneeded side-cases are omitted for
// simplicity.  See:
//     https://en.wikipedia.org/wiki/Elliptic_curve_point_multiplication#Point_addition
//     https://stackoverflow.com/a/31089415
//     https://en.wikipedia.org/wiki/Modular_multiplicative_inverse#Using_Euler's_theorem
//     https://crypto.stanford.edu/pbc/notes/elliptic/explicit.html
func (R *Point) ec_point_add(P, Q *Point) *Point {
    //if P.x == nil {
    //    R.x = new(big.Int)
    //    R.y = new(big.Int)
    //    R.x.Set(Q.x)
    //    R.y.Set(Q.y)
    //    return R
    //}

    R.x.Add(P.x, big.NewInt(1))
    R.y.Add(P.y, big.NewInt(1))
    return R
}

func main() {
    fmt.Println("Example from Mastering Bitcoin, pages 69-70.")

    // A private key must be a whole number from 1 to
    //     0xfffffffffffffffffffffffffffffffebaaedce6af48a03bbfd25e8cd0364140,
    // one less than the order of the base point (or "generator point")
    // G.  See:
    //     https://en.bitcoin.it/wiki/Private_key//Range_of_valid_ECDSA_private_keys
    var private_key = new(big.Int)
    private_key.SetString(
        "038109007313a5807b2eccc082c8c3fbb988a973cacf1a7df9ce725c31b14776", 16)
    fmt.Println("\tprivate_key:")
    fmt.Println(private_key)

    // Secp256k1 parameters.  See:
    //     https://en.bitcoin.it/wiki/Secp256k1
    //     https://www.secg.org/sec2-v2.pdf - Section 2.4.1.
    var p = big.NewInt(2)
    p.Exp(p, big.NewInt(256), nil)
    p.Sub(p, big.NewInt(1<<32 + 1<<9 + 1<<8 + 1<<7 + 1<<6 + 1<<4 + 1))
    //fmt.Println(p)

    var G = Point{new(big.Int), new(big.Int)}
    //G.x = new(big.Int)
    G.x.SetString(
        "79be667ef9dcbbac55a06295ce870b07029bfcdb2dce28d959f2815b16f81798", 16)
    //G.y = new(big.Int)
    G.y.SetString(
        "483ada7726a3c4655da4fbfc0e1108a8fd17b448a68554199c47d08ffb10d4b8", 16)
    fmt.Println(G.x)
    fmt.Println(G.y)

    var A = Point{new(big.Int), new(big.Int)}
    G.ec_point_add(&G, &G)
    fmt.Println(A.x)
    fmt.Println(A.y)
    fmt.Println(G.y)
}
