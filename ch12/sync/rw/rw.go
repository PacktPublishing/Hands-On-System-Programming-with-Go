package main

import (
	"fmt"
	"sync"
	"time"
)

type list struct {
	m     sync.RWMutex
	value []string
}

func (l *list) contains(v string) bool {
	for _, s := range l.value {
		if s == v {
			return true
		}
	}
	return false
}

func (l *list) Len() int {
	l.m.RLock()
	v := len(l.value)
	l.m.RUnlock()
	return v
}

func (l *list) Contains(v string) bool {
	l.m.RLock()
	found := l.contains(v)
	l.m.RUnlock()
	return found
}

func (l *list) Add(v string) bool {
	l.m.Lock()
	defer l.m.Unlock()
	if l.contains(v) {
		return false
	}
	l.value = append(l.value, v)
	return true
}

func main() {
	var src = []string{
		"Ryu", "Ken", "E. Honda", "Guile",
		"Chun-Li", "Blanka", "Zangief", "Dhalsim",
	}
	var l list

	for i := 0; i < 10; i++ {
		go func(i int) {
			for _, s := range src {
				go func(s string) {
					if !l.Contains(s) {
						if l.Add(s) {
							fmt.Println(i, "add", s)
						} else {
							fmt.Println(i, "too slow", s)
						}
					}
				}(s)
			}
		}(i)
	}

	time.Sleep(500 * time.Millisecond)
}
