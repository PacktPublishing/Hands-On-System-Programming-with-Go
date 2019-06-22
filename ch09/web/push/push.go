package main

import (
	"fmt"
	"net/http"
)

func main() {
	const imgPath = "/image.svg"
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		pusher, ok := w.(http.Pusher)
		if ok {
			fmt.Println("Push /image")
			pusher.Push(imgPath, nil)
		}
		w.Header().Add("Content-Type", "text/html")
		fmt.Fprintf(w, `<html><body><img src="%s"/>`+
			"</body></html>\n", imgPath)
	})
	http.HandleFunc(imgPath, func(w http.ResponseWriter, r *http.Request) {
		w.Header().Add("Content-Type", "image/svg+xml")
		fmt.Fprint(w, `<?xml version="1.0" standalone="no"?>
<svg xmlns="http://www.w3.org/2000/svg">
  <rect width="150" height="150" style="fill:blue"/>
</svg>`)
	})
	if err := http.ListenAndServeTLS(":3000", "server.crt", "server.key", nil); err != nil {
		fmt.Println(err)
	}
}
