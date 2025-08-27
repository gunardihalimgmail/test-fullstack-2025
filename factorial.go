package main

import (
	"fmt"
	"math"
)

func funcFactorial(n int) uint64 {
	if n == 0 {
		return 1
	}
	result := uint64(1)
	for i := 2; i <= n; i++ {
		result *= uint64(i)
	}
	return result
}

func funcHitung(n int) uint64 {

	factorial := float64(funcFactorial(n))
	pangkat := math.Pow(2, float64(n))
	hasil := factorial / pangkat

	return uint64(math.Ceil(hasil))
}

func main() {
	for i := 0; i <= 10; i++ {
		fmt.Printf("funcHitung(%d) = %d\n", i, funcHitung(i))
	}
}
