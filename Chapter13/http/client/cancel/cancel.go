package main

import (
	"context"
	"log"
	"math/rand"
	"net/http"
	"sync"
	"time"
)

func main() {
	const addr = "localhost:8080"
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		d := time.Second * time.Duration(rand.Intn(10))
		log.Println("wait", d)
		time.Sleep(d)
	})
	go func() {
		if err := http.ListenAndServe(addr, nil); err != nil {
			log.Fatalln(err)
		}
	}()

	ctx, canc := context.WithCancel(context.Background())
	ch, o, wg := make(chan int), sync.Once{}, sync.WaitGroup{}
	wg.Add(10)
	for i := 0; i < 10; i++ {
		go func(i int) {
			defer wg.Done()
			req, _ := http.NewRequest(http.MethodGet, "http://"+addr, nil)
			if _, err := http.DefaultClient.Do(req.WithContext(ctx)); err != nil {
				log.Println(i, err)
				return
			}
			o.Do(func() { ch <- i })
		}(i)
	}
	log.Println("received", <-ch)
	canc()
	log.Println("cancelling")
	wg.Wait()
}
