package main

import (
	"fmt"
	"math/big"
)

func main() {
	fmt.Println("highest prime number: ", highestPrimeNumber(100))

}

func highestPrimeNumber(num int64) int64 {
Prime:
	for {
		switch {
		case big.NewInt(num).ProbablyPrime(0):
			break Prime
		default:
			num--
		}
	}
	return num
}
