package main

import (
	"fmt"
	"math/big"
)

func b58(data []byte) string {
	base := big.NewInt(58)                                                        // set new big int for dividing data bytes with 58
	remainder := new(big.Int)                                                     // to store remainder of the data bytes
	const alphabet = "123456789ABCDEFGHJKLMNPQRSTUVWXYZabcdefghijkmnopqrstuvwxyz" // base 58 character set

	//fmt.Println("Length of base58 characters: ", len(alphabet))
	fmt.Println("Recieved data bytes: ", data)

	x := new(big.Int).SetBytes(data) // convert bytes to big integer
	//fmt.Println("Bytes to big integer: ", x)
	var output_string []byte // To store base58 coverted value to a temp variable

	for x.Cmp(big.NewInt(0)) != 0 {
		x.DivMod(x, base, remainder)
		//fmt.Println("X", x, "\nRemainder", remainder)
		fmt.Printf("alphabet[%d]: %s\n", remainder.Int64(), string(alphabet[remainder.Int64()]))
		output_string = append(output_string, alphabet[remainder.Int64()]) // Appending the value stored at position alphabet[remainder] to output_string
	}
	//fmt.Println("output_string value after first loop: ", string(output_string))

	for i := 0; data[i] == 0; i++ {
		output_string = append(output_string, alphabet[0])
	}
	fmt.Println("output_string value after second loop: ", string(output_string))

	// Reverse the array.
	for i, j := 0, len(output_string)-1; i < j; i, j = i+1, j-1 {
		output_string[i], output_string[j] = output_string[j], output_string[i]
	}
	return string(output_string)
}

func b58decode(data string) []byte {
	fmt.Println(data)
	dataByte := []byte(data)
	return dataByte
}

func main() {
	dog := b58([]byte("cat"))
	fmt.Println("base58 encode:", dog, "\n")

	decoded := b58decode(dog)
	fmt.Println("base58 decode:", decoded)
}
