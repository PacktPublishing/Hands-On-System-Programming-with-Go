package main

import (
    "context"
    "fmt"
    "log"
    "math/rand"
    "time"
)

type msg struct {
    value string
    done  chan struct{}
}

func (m *msg) Wait() {
    <-m.done
}

func (m *msg) Done() {
    m.done <- struct{}{}
}

func send(ctx context.Context, v string) <-chan msg {
    ch := make(chan msg)
    go func() {
        done := make(chan struct{})
        for i := 0; ; i++ {
            time.Sleep(time.Duration(float64(time.Second/2) * rand.Float64()))
            m := msg{fmt.Sprintf("%s msg-%d", v, i), done}
            select {
            case <-ctx.Done():
                close(ch)
                return
            case ch <- m:
                m.Wait()
            }
        }
    }()
    return ch
}

func merge(ctx context.Context, sources ...<-chan msg) <-chan msg {
    ch := make(chan msg)
    go func() {
        <-ctx.Done()
        close(ch)
    }()
    for i := range sources {
        go func(i int) {
            for {
                select {
                case v := <-sources[i]:
                    select {
                    case <-ctx.Done():
                        return
                    case ch <- v:
                    }
                }
            }
        }(i)
    }
    return ch
}

func main() {
    ctx, canc := context.WithTimeout(context.Background(), time.Second)
    defer canc()
    sources := make([]<-chan msg, 5)
    for i := range sources {
        sources[i] = send(ctx, fmt.Sprint("src-", i))
    }
    msgs := make([]msg, 0, len(sources))
    start := time.Now()
    for v := range merge(ctx, sources...) {
        msgs = append(msgs, v)
        log.Println(v.value, time.Since(start))
        if len(msgs) == len(sources) {
            log.Println("*** done ***")
            for _, m := range msgs {
                m.Done()
            }
            msgs = msgs[:0]
            start = time.Now()
        }
    }
}
