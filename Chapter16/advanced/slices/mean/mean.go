package main

import "fmt"

/*
double mean(int len, double *a) {
    if (a == NULL || len == 0) {
        return 0;
    }
    double m = 0;
    for (int i = 0; i < len; i++) {
        m+=a[i];
    }
    return m / len;
}
*/
import "C"

func mean(a []float64) float64 {
    if len(a) == 0 {
        return 0
    }
    return float64(C.mean(C.int(len(a)), (*C.double)(&a[0])))
}

func mean2(a []float64) float64 {
    if len(a) == 0 {
        return 0
    }
    return float64(C.mean(C.int(len(a)*2), (*C.double)(&a[0])))
}

var a = make([]float64, 10)

func init() {
    for i := range a {
        a[i] = float64(i + 1)
    }
}

func main() {
    cases := [][]float64{a, a[1:4], a[:0], nil}
    for _, slice := range cases {
        fmt.Println(slice, mean(slice))
    }
    for _, slice := range cases {
        fmt.Println(slice, mean2(slice))
    }
}
