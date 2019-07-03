package main

import (
	"fmt"
)

func main() {
	done := make(chan struct{}, 10000)
	var a = 0
	for i := 0; i < cap(done); i++ {
		go func(i int) {
			if i%2 == 0 {
				a++
			} else {
				a--
			}
			done <- struct{}{}
		}(i)
	}
	for i := 0; i < cap(done); i++ {
		<-done
	}
	fmt.Println(a)
}
