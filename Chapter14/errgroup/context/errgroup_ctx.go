package main

import (
    "context"
    "errors"
    "fmt"
    "log"
    "math/rand"
    "time"

    "golang.org/x/sync/errgroup"
)

func sender(ctx context.Context, ch chan<- string, n int) func() error {
    return func() (err error) {
        for i := 0; ; i++ {
            if rand.Intn(100) == 42 {
                return errors.New("the answer")
            }
            select {
            case ch <- fmt.Sprintf("[%d]%d", n, i):
            case <-ctx.Done():
                return nil
            }
        }
    }
}

func main() {
    eg, ctx := errgroup.WithContext(context.Background())
    ch := make(chan string)
    for i := 0; i < 10; i++ {
        eg.Go(sender(ctx, ch, i))
    }
    go func() {
        for s := range ch {
            log.Println(s)
        }
    }()
    if err := eg.Wait(); err != nil {
        log.Println("Error:", err)
    }
    close(ch)
    log.Println("waiting...")
    time.Sleep(time.Second)
}
