package main

import (
	"fmt"
	"sync"
	"sync/atomic"
)

type clicker int32

func (c *clicker) Click() int32 {
	return atomic.AddInt32((*int32)(c), 1)
}

func (c *clicker) Reset() {
	atomic.StoreInt32((*int32)(c), 0)
}

func (c *clicker) Value() int32 {
	return atomic.LoadInt32((*int32)(c))
}

func main() {
	c := clicker(0)
	wg := sync.WaitGroup{}
	// 2*iteration + reset at 5
	wg.Add(21)
	for i := 0; i < 10; i++ {
		go func() {
			c.Click()
			fmt.Println("click")
			wg.Done()
		}()
		go func() {
			fmt.Println("load", c.Value())
			wg.Done()
		}()
		if i == 0 || i%5 != 0 {
			continue
		}
		fmt.Println(i)
		go func() {
			c.Reset()
			fmt.Println("reset")
			wg.Done()
		}()
	}
	wg.Wait()
}
