package main

import (
	"context"
	"fmt"
	"sync"
	"time"

	"golang.org/x/sync/semaphore"
)

func main() {
	s := semaphore.NewWeighted(int64(5))
	ctx, canc := context.WithTimeout(context.Background(), time.Second)
	defer canc()
	wg := sync.WaitGroup{}
	wg.Add(20)
	for i := 0; i < 20; i++ {
		go func(i int) {
			defer wg.Done()
			if err := s.Acquire(ctx, 1); err != nil {
				fmt.Println(i, err)
				return
			}
			go func(i int) {
				fmt.Println(i)
				time.Sleep(time.Second / 2)
				s.Release(1)
			}(i)
		}(i)
	}
	wg.Wait()
}
