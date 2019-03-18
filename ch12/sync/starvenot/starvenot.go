package main

import (
	"fmt"
	"sync"
	"time"
)

type counter struct {
	m     sync.Mutex
	value int
}

func (c *counter) Write(i int) {
	c.m.Lock()
	time.Sleep(time.Millisecond * 100)
	c.value = i
	c.m.Unlock()
}

func (c *counter) Value() int {
	c.m.Lock()
	time.Sleep(time.Millisecond * 100)
	a := c.value
	c.m.Unlock()
	return a
}

func main() {
	var c counter
	t1 := time.NewTicker(time.Millisecond * 50)
	time.AfterFunc(time.Second*2, t1.Stop)
	for {
		select {
		case <-t1.C:
			go func() {
				t := time.Now()
				c.Value()
				fmt.Println("val", time.Since(t))
			}()
			go func() {
				t := time.Now()
				c.Write(0)
				fmt.Println("inc", time.Since(t))
			}()
		case <-time.After(time.Millisecond * 200):
			return
		}
	}
}
