package main

import (
    "log"
    "net/http"
    "time"

    "golang.org/x/sync/errgroup"
)

func visitor(url string) func() error {
    return func() (err error) {
        s := time.Now()
        defer func() {
            log.Println(url, time.Since(s), err)
        }()
        var resp *http.Response
        if resp, err = http.Get(url); err != nil {
            return
        }
        return resp.Body.Close()
    }
}

func main() {
    eg := errgroup.Group{}
    var urlList = []string{
        "http://www.golang.org/",
        "http://invalidwebsite.hey/",
        "http://www.google.com/",
    }
    for _, url := range urlList {
        eg.Go(visitor(url))
    }
    if err := eg.Wait(); err != nil {
        log.Fatalln("Error:", err)
    }
}
