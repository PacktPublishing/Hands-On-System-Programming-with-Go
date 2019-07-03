package main

import (
	"context"
	"log"
	"net/http"
	"time"
)

func main() {
	const addr = "localhost:8080"
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		time.Sleep(time.Second * 5)
	})
	go func() {
		if err := http.ListenAndServe(addr, nil); err != nil {
			log.Fatalln(err)
		}
	}()
	req, _ := http.NewRequest(http.MethodGet, "http://"+addr, nil)
	ctx, canc := context.WithTimeout(context.Background(), time.Second*2)
	defer canc()
	time.Sleep(time.Second)
	if _, err := http.DefaultClient.Do(req.WithContext(ctx)); err != nil {
		log.Fatalln(err)
	}

}
