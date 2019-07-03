package main

import (
	"fmt"
	"sync"
)

func main() {
	var m sync.Map
	fmt.Println(m.LoadOrStore(1, 2))
	fmt.Println(m.Load(1))
	fmt.Println(m.LoadOrStore(1, 3))
	fmt.Println(m.Load(1))
	fmt.Println(m.LoadOrStore(1, 4))
	fmt.Println(m.Load(1))
}
