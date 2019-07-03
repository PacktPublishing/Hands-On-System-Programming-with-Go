package main

/*
#include <stdio.h>
#include <stdlib.h>
#include <string.h>

char* reverseString(char* s) {
    int l = strlen(s);
    for (int i=0; i < l/2; i++) {
        char a = s[i];
        s[i] = s[l-1-i];
        s[l-1-i] = a;
    }
    return s;
}
*/
import "C"

import (
    "fmt"
    "unsafe"
)

func main() {
    b1 := []byte("A byte slice")
    c1  := C.CBytes(b1)
    fmt.Printf("Go ptr: %p\n", b1)
    fmt.Printf("C ptr:  %p\n", c1)
    defer C.free(c1)
    c2 := unsafe.Pointer(C.reverseString((*C.char)(c1)))
    b2 := C.GoBytes(c2, C.int(len(b1)))
    fmt.Printf("Go ptr: %p\n", b2)
    fmt.Printf("%q -> %q", b1, b2)
}