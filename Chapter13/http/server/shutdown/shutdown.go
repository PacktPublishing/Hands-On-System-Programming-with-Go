package main

import (
	"context"
	"log"
	"net/http"
	"time"
)

func main() {
	mux := http.NewServeMux()
	server := http.Server{
		Addr:    ":3000",
		Handler: mux,
	}
	ctx, canc := context.WithCancel(context.Background())
	defer canc()
	mux.HandleFunc("/shutdown", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("OK"))
		canc()
	})
	go func() {
		if err := server.ListenAndServe(); err != nil {
			if err != http.ErrServerClosed {
				log.Fatal(err)
			}
		}
	}()
	select {
	case <-ctx.Done():
		ctx, canc := context.WithTimeout(context.Background(), time.Second*5)
		defer canc()
		if err := server.Shutdown(ctx); err != nil {
			log.Fatalln("Shutdown:", err)
		} else {
			log.Println("Shutdown:", "ok")
		}
	}
}
