package main

/*
#include <stdio.h>
#include <stdlib.h>
#include <string.h>

void reverseString(char* s) {
    int l = strlen(s);
    for (int i=0; i < l/2; i++) {
        char a = s[i];
        s[i] = s[l-1-i];
        s[l-1-i] = a;
    }
}
*/
import "C"

import (
	"fmt"
	"unsafe"
)

func main() {
	b1 := []byte("A byte slice")
	fmt.Printf("Slice: %q\n", b1)
	C.reverseString((*C.char)(unsafe.Pointer(&b1[0])))
	fmt.Printf("Slice: %q\n", b1)
}
