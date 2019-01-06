package main

import (
	"fmt"
	"io/ioutil"
	"os"
)

func main() {
	if len(os.Args) != 3 {
		fmt.Println("Please specify a path and some content")
		return
	}
	// the second argument, the content, needs to be casted to a byte slice
	if err := ioutil.WriteFile(os.Args[1], []byte(os.Args[2]), 0644); err != nil {
		fmt.Println("Error:", err)
	}
}
