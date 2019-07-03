package main

import "C"

import (
	"fmt"
	"math"
)

func main() {
	a := float64(math.Pi)

	fmt.Println(a)
	fmt.Println(C.float(a))
	fmt.Println(C.double(a))
	fmt.Println(C.double(C.float(a)) - C.double(a))
}
