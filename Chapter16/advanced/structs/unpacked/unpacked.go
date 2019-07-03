package main

/*
typedef struct{
    unsigned char a;
    char b;
    int c;
    unsigned int d;
    char e[10];
} myStruct;
*/
import "C"
import "log"

func main() {
	v := C.myStruct{
		a: C.uchar('A'),
		b: C.char('Z'),
		c: C.int(100),
		d: C.uint(10),
		e: [10]C.char{'h', 'e', 'l', 'l', 'o'},
	}
	log.Printf("%#v", v)
}
