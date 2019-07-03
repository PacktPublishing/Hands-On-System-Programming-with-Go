package main

import (
	"bufio"
	"context"
	"flag"
	"fmt"
	"io/ioutil"
	"os"
	"os/signal"
	"path/filepath"
	"strings"
	"sync"
	"syscall"
)

var o = Options{}

func init() {
	var exclude string
	flag.BoolVar(&o.Contents, "c", false, "Search contents")
	flag.StringVar(&exclude, "x", "", "Exclude folders (use : as separator)")
	flag.Parse()
	if exclude != "" {
		o.Exclude = strings.Split(exclude, ":")
	}
}

func main() {
	args := flag.Args()
	if len(args) != 2 {
		fmt.Printf("Usage `%s <file> <search_term>`\n", filepath.Base(os.Args[0]))
		os.Exit(2)
	}

	ctx, cancel := context.WithCancel(context.Background())
	c := make(chan os.Signal)
	signal.Notify(c, syscall.SIGTERM, syscall.SIGINT)
	go func() {
		<-c
		cancel()
	}()
	for r := range FileSearch(ctx, args[0], args[1], &o) {
		if r.Err != nil {
			fmt.Printf("%s - error: %s\n", r.File, r.Err)
			continue
		}
		if !o.Contents {
			fmt.Printf("%s - match\n", r.File)
			continue
		}
		fmt.Printf("%s - matches:\n", r.File)
		for _, m := range r.Matches {
			fmt.Printf("\t%d:%s\n", m.Line, m.Text)
		}
	}

}

type Result struct {
	Err     error
	File    string
	Matches []Match
}

type Match struct {
	Line int
	Text string
}

type Options struct {
	Contents bool
	Exclude  []string
}

func FileSearch(ctx context.Context, root, term string, o *Options) <-chan Result {
	ch, wg, once := make(chan Result), sync.WaitGroup{}, sync.Once{}
	wg.Add(1)
	go fileSearch(ctx, ch, &wg, root, term, o)
	go func() {
		wg.Wait()
		fmt.Println("* Search done *")
		once.Do(func() {
			close(ch)
		})
	}()
	go func() {
		<-ctx.Done()
		fmt.Println("* Context done *")
		once.Do(func() {
			close(ch)
		})
	}()
	return ch
}

func fileSearch(ctx context.Context, ch chan<- Result, wg *sync.WaitGroup, file, term string, o *Options) {
	defer wg.Done()
	_, name := filepath.Split(file)
	if o != nil {
		for _, e := range o.Exclude {
			if e == name {
				return
			}
		}
	}
	info, err := os.Stat(file)
	if err != nil {
		select {
		case <-ctx.Done():
			return
		default:
			ch <- Result{File: file, Err: err}
		}
		return
	}
	if info.IsDir() {
		files, err := ioutil.ReadDir(file)
		if err != nil {
			select {
			case <-ctx.Done():
			default:
				ch <- Result{File: file, Err: err}
			}
			return
		}
		select {
		case <-ctx.Done():
			return
		default:
			wg.Add(len(files))
			for _, f := range files {
				go fileSearch(ctx, ch, wg, filepath.Join(file, f.Name()), term, o)
			}
		}
		return
	}
	if o == nil || !o.Contents {
		if name == term {
			select {
			case <-ctx.Done():
			default:
				ch <- Result{File: file}
			}
		}
		return
	}
	//file search
	f, err := os.Open(file)
	if err != nil {
		select {
		case <-ctx.Done():
		default:
			ch <- Result{File: file, Err: err}
		}
		return
	}
	defer f.Close()

	scanner, matches, line := bufio.NewScanner(f), []Match{}, 1
	for scanner.Scan() {
		select {
		case <-ctx.Done():
			break
		default:
			if text := scanner.Text(); strings.Contains(text, term) {
				matches = append(matches, Match{Line: line, Text: text})
			}
			line++
		}
	}

	select {
	case <-ctx.Done():
		break
	default:
		if err := scanner.Err(); err != nil {
			ch <- Result{File: file, Err: err}
			return
		}
		if len(matches) != 0 {
			ch <- Result{File: file, Matches: matches}
		}
	}
}
