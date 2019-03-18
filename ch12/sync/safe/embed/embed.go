package main

import (
	"fmt"
	"sync"
)

type counter struct {
	m     sync.Mutex
	value int
}

func (c *counter) Incr() {
	c.m.Lock()
	c.value++
	c.m.Unlock()
}

func (c *counter) Decr() {
	c.m.Lock()
	c.value--
	c.m.Unlock()
}

func (c *counter) Value() int {
	c.m.Lock()
	a := c.value
	c.m.Unlock()
	return a
}
func main() {
	done := make(chan struct{}, 10000)
	var a = counter{}
	for i := 0; i < cap(done); i++ {
		go func(i int) {
			if i%2 == 0 {
				a.Incr()
			} else {
				a.Decr()
			}
			done <- struct{}{}
		}(i)
	}
	for i := 0; i < cap(done); i++ {
		<-done
	}
	fmt.Println(a.Value())
}
