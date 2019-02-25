package main

import (
	"fmt"
	"log"
	"net/http"
)

type customHandler int

func (c *customHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "%d", *c)
	*c++
}

func main() {
	mux := http.NewServeMux()
	mux.HandleFunc("/hello", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "Hello!")
	})
	mux.HandleFunc("/bye", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "Goodbye!")
	})
	mux.HandleFunc("/error", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(w, "An error occurred!")
	})
	mux.Handle("/custom", new(customHandler))
	if err := http.ListenAndServe(":3000", mux); err != nil {
		log.Fatal(err)
	}
}
