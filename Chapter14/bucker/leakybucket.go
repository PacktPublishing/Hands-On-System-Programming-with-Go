package main

import (
    "context"
    "fmt"
    "log"
    "math/rand"
    "sync/atomic"
    "time"
)

type bucket struct {
    capacity uint64
    status   uint64
}

func newBucket(ctx context.Context, cap uint64, rate time.Duration) *bucket {
    b := bucket{capacity: cap, status: cap}
    go func() {
        t := time.NewTicker(rate)
        for {
            select {
            case <-t.C:
                fmt.Println("")
                atomic.StoreUint64(&b.status, b.capacity)
            case <-ctx.Done():
                t.Stop()
                return
            }
        }
    }()
    return &b
}

func (b *bucket) Add(n uint64) uint64 {
    for {
        r := atomic.LoadUint64(&b.status)
        if r == 0 {
            return 0
        }
        if n > r {
            n = r
        }
        if !atomic.CompareAndSwapUint64(&b.status, r, r-n) {
            continue
        }
        return n
    }
}

func main() {
    ctx, canc := context.WithTimeout(context.Background(), time.Second)
    defer canc()
    start := time.Now()
    b := newBucket(ctx, 10, time.Second/5)
    t := time.Second / 10
    for i := 0; i < 5; i++ {
        c := client{
            name:  fmt.Sprint(i),
            b:     b,
            sleep: t,
            max:   5,
        }
        go c.Run(ctx, start)
    }
    <-ctx.Done()
}

type client struct {
    name  string
    max   int
    b     *bucket
    sleep time.Duration
}

func (c client) Run(ctx context.Context, start time.Time) {
    for {
        select {
        case <-ctx.Done():
            return
        default:
            n := 1 + rand.Intn(c.max-1)
            time.Sleep(c.sleep)
            e := time.Since(start).Seconds()
            a := c.b.Add(uint64(n))
            log.Printf("%s tries to take %d after %.02fs, takes %d", c.name, n, e, a)
        }
    }
}
