package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"sync/atomic"
)

type keyType struct{}

var key = &keyType{}

var counter int32

func WithKey(ctx context.Context) context.Context {
	return context.WithValue(ctx, key, atomic.AddInt32(&counter, 1))
}

func GetKey(ctx context.Context) (int32, bool) {
	v := ctx.Value(key)
	if v == nil {
		return 0, false
	}
	return v.(int32), true
}

func AssignKeyHandler(h http.Handler) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		if ctx == nil {
			ctx = context.Background()
		}
		if _, ok := GetKey(ctx); !ok {
			ctx = WithKey(ctx)
		}
		h.ServeHTTP(w, r.WithContext(ctx))
	}
}

func ReadFileHandler(root string) http.HandlerFunc {
	root = filepath.Clean(root)
	return func(w http.ResponseWriter, r *http.Request) {
		k, _ := GetKey(r.Context())
		path := filepath.Join(root, r.URL.Path)
		log.Printf("[%d] requesting path %s", k, path)
		if !strings.HasPrefix(path, root) {
			http.Error(w, "not found", http.StatusNotFound)
			log.Printf("[%d] unauthorized %s", k, path)
			return
		}
		if stat, err := os.Stat(path); err != nil || stat.IsDir() {
			http.Error(w, "not found", http.StatusNotFound)
			log.Printf("[%d] not found %s", k, path)
			return
		}
		http.ServeFile(w, r, path)
		log.Printf("[%d] ok: %s", k, path)
	}
}

func main() {
	home, err := os.UserHomeDir()
	if err != nil {
		log.Fatal(err)
	}
	tmp := os.TempDir()
	mux := http.NewServeMux()
	server := http.Server{
		Addr:    ":3000",
		Handler: mux,
	}

	mux.Handle("/tmp/", http.StripPrefix("/tmp/", AssignKeyHandler(ReadFileHandler(tmp))))
	mux.Handle("/home/", http.StripPrefix("/home/", AssignKeyHandler(ReadFileHandler(home))))
	if err := server.ListenAndServe(); err != nil {
		if err != http.ErrServerClosed {
			log.Fatal(err)
		}
	}
}
