package main

/*
#include <stdio.h>
#include <stdlib.h>

void customPrint(char* s) {
	printf("%s\n", s);
}
*/
import "C"

import "unsafe"

func main() {
	s := C.CString(`Printing to stdout with CGO
Using <stdio.h> and <stdlib.h>`)
	defer C.free(unsafe.Pointer(s))
	C.customPrint(s)
}