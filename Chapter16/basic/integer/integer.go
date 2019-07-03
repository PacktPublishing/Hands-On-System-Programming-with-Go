package main

import "C"

import "fmt"

func main() {
	a := int64(0x1122334455667788)

	// a fits in 64 bits
	fmt.Println(a)
	// short overflows, it's 16
	fmt.Println(C.short(a), int16(0x7788))
	// long also overflows, it's 32
	fmt.Println(C.long(a), int32(0x55667788))
	// longlong is okay, it's 64
	fmt.Println(C.longlong(a), int64(0x1122334455667788))
}
