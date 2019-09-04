package test

import (
	"fmt"
	"testing"
)

func TestSqrt(t *testing.T) {
	z := Sqrt(5)
	fmt.Println("Hello ", z)
}

func Sqrt(x float64) float64 {
	z := 0.0
	for i := 0; i < 1000; i++ {
		z -= (z*z - x) / (2 * x)
	}
	return z
}
