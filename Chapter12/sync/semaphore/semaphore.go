package main

import (
    "context"
    "fmt"
    "log"
    "time"

    "golang.org/x/sync/semaphore"
)

func main() {
    s := semaphore.NewWeighted(int64(10))
    ctx := context.Background()
    for i := 0; i < 20; i++ {
        if err := s.Acquire(ctx, 1); err != nil {
            log.Fatal(err)
        }
        go func(i int) {
            fmt.Println(i)
            s.Release(1)
        }(i)
    }
    time.Sleep(time.Second)
}
