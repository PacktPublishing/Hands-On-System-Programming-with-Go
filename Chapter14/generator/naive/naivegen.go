package main

import (
	"fmt"
	"time"
)

type genInt64 int64

func (g *genInt64) Next() int64 {
	*g++
	return int64(*g)
}

func main() {
	var g genInt64
	for i := 0; i < 1000; i++ {
		go func(i int) {
			fmt.Println(i, g.Next())
		}(i)
	}
	time.Sleep(time.Second)
}
