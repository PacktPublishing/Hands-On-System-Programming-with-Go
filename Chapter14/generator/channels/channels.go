package main

import "fmt"

func GenInt64() chan int64 {
    ch := make(chan int64)
    go func() {
        for i := int64(0); ; i++ {
            ch <- i
        }
    }()
    return ch
}

func main() {
    ch1, ch2 := GenInt64(), GenInt64()
    for i := 0; i < 20; i++ {
        select {
        case v := <-ch1:
            fmt.Println("ch 1", v)
        case v := <-ch2:
            fmt.Println("ch 2", v)
        }
    }
}
