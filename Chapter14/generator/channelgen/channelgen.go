package main

import (
	"context"
	"fmt"
	"time"
)

type genInt64 struct {
	ch chan int64
}

func (g genInt64) Next() int64 {
	return <-g.ch
}

func NewGenInt64(ctx context.Context) genInt64 {
	g := genInt64{ch: make(chan int64)}
	go func() {
		for i := int64(0); ; i++ {
			select {
			case g.ch <- i:
				// do nothing
			case <-ctx.Done():
				close(g.ch)
				return
			}
		}
	}()
	return g
}

func main() {
	g := NewGenInt64(context.Background())
	for i := 0; i < 1000; i++ {
		go func(i int) {
			fmt.Println(i, g.Next())
		}(i)
	}
	time.Sleep(time.Second)
}
