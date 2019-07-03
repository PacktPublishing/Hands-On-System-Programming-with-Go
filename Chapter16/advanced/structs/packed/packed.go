package main

/*
#include "stdio.h"
#pragma pack(1)
typedef struct{
    unsigned char a;
    char b;
    int c;
    unsigned int d;
    char e[10];
}myStruct;

myStruct makeStruct(){
    myStruct p;
    p.a = 'A';
    p.b = 'Z';
    p.c = 100;
    p.d = 10;
    p.e[0] = 'h';
    p.e[1] = 'e';
    p.e[2] = 'l';
    p.e[3] = 'l';
    p.e[4] = 'o';
    p.e[5] = '\0';
    p.e[6] = '\0';
    p.e[7] = '\0';
    p.e[8] = '\0';
    p.e[9] = '\0';
    return p;
}

*/
import "C"
import (
	"bytes"
	"encoding/binary"
	"log"
	"unsafe"
)

type myStruct struct {
	a uint8
	b int8
	c int32
	d uint32
	e [10]uint8
}

func unpack(i *C.myStruct) (m myStruct) {
	b := bytes.NewBuffer(C.GoBytes(unsafe.Pointer(i), C.sizeof_myStruct))
	for _, v := range []interface{}{&m.a, &m.b, &m.c, &m.d, &m.e} {
		binary.Read(b, binary.LittleEndian, v)
	}
	return
}

func main() {
	v := C.myStruct{
		a: C.uchar('A'),
		b: C.char('Z'),
		//c: C.int(100),
		//d: C.uint(10),
		e: [10]C.char{'h', 'e', 'l', 'l', 'o'},
	}
	log.Printf("%#v", v)

	v = C.makeStruct()
	log.Printf("%#v", v)
	m := unpack(&v)
	log.Printf("%#v", m)
}
