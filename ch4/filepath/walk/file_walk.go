package main

import (
	"fmt"
	"os"
	"path/filepath"
)

func main() {
	if len(os.Args) != 2 { // ensure path is specified
		fmt.Println("Please specify a path.")
		return
	}
	root, err := filepath.Abs(os.Args[1]) // get absolute path
	if err != nil {
		fmt.Println("Cannot get absolute path:", err)
		return
	}
	fmt.Println("Listing files in", root)
	var c struct {
		files int
		dirs  int
	}
	filepath.Walk(root, func(path string, info os.FileInfo, err error) error {
		// walk the tree to count files and folders
		if info.IsDir() {
			c.dirs++
		} else {
			c.files++
		}
		fmt.Println("-", path)
		return nil
	})
	fmt.Printf("Total: %d files in %d directories", c.files, c.dirs)
}
