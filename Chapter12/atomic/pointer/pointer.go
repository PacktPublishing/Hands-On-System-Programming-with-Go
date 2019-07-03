package main

import (
	"log"
	"sync/atomic"
	"unsafe"
)

func main() {
	v1, v2 := 10, 100
	p1, p2 := &v1, &v2
	log.Printf("P1: %v, P2: %v", *p1, *p2)
	atomic.SwapPointer((*unsafe.Pointer)(unsafe.Pointer(&p1)), unsafe.Pointer(p2))
	log.Printf("P1: %v, P2: %v", *p1, *p2)
	v1 = -10
	log.Printf("P1: %v, P2: %v", *p1, *p2)
	v2 = 3
	log.Printf("P1: %v, P2: %v", *p1, *p2)
}
