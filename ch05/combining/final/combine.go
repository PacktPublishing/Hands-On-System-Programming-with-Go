package main

import (
	"bytes"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
)

type queryWriter struct {
	Query []byte
	io.Writer
}

func (q queryWriter) Write(b []byte) (n int, err error) {
	lines := bytes.Split(b, []byte{'\n'})
	l := len(q.Query)
	for _, b := range lines {
		i := bytes.Index(b, q.Query)
		if i == -1 {
			continue
		}
		for _, s := range [][]byte{
			b[:i],              // what's before the match
			[]byte("\x1b[31m"), //star red color
			b[i : i+l],         // match
			[]byte("\x1b[39m"), // default color
			b[i+l:],            // whatever is left
		} {
			v, err := q.Writer.Write(s)
			n += v
			if err != nil {
				return 0, err
			}
		}
		fmt.Fprintln(q.Writer)
	}
	return len(b), nil
}

func main() {
	if len(os.Args) < 3 {
		fmt.Println("Please specify a path and a search string.")
		return
	}
	root, err := filepath.Abs(os.Args[1]) // get absolute path
	if err != nil {
		fmt.Println("Cannot get absolute path:", err)
		return
	}
	query := []byte(strings.Join(os.Args[2:], " "))
	fmt.Printf("Searching for %q in %s...\n", query, root)
	err = filepath.Walk(root, func(path string, info os.FileInfo, err error) error {
		if info.IsDir() {
			return nil
		}
		fmt.Println(path)
		f, err := os.Open(path)
		if err != nil {
			return err
		}
		defer f.Close()

		_, err = ioutil.ReadAll(io.TeeReader(f, queryWriter{Query: query, Writer: os.Stdout}))
		return err
	})
	if err != nil {
		fmt.Println(err)
	}

}
