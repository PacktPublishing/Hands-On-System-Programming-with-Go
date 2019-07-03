package main

import (
	"context"
	"fmt"
	"reflect"
	"time"
)

func main() {
	testMap()
	testStruct()
}

func baseTest(fn1, fn2 func(int)) {
	ctx, canc := context.WithTimeout(context.Background(), time.Second)
	defer canc()
	go func() {
		for i := 0; ; i++ {
			select {
			case <-ctx.Done():
				return
			default:
				fn1(i)
			}
		}
	}()
	go func() {
		for i := 0; ; i++ {
			select {
			case <-ctx.Done():
				return
			default:
				fn2(i)
			}
		}
	}()
	<-ctx.Done()
}

func testMap() {
	m1, m2 := make(map[int]int), make(map[int]int)
	m := reflect.ValueOf(m2)
	baseTest(func(i int) { m1[i] = i }, func(i int) {
		v := reflect.ValueOf(i)
		m.SetMapIndex(v, v)
	})
	fmt.Printf("normal	%d\n", len(m1))
	fmt.Printf("reflect	%d\n", len(m2))
}

func testStruct() {
	type T struct {
		Field int
	}
	var m1, m2 T
	m := reflect.ValueOf(&m2).Elem()
	baseTest(func(i int) { m1.Field++ }, func(i int) {
		f := m.Field(0)
		f.SetInt(int64(f.Interface().(int) + 1))
	})
	fmt.Printf("normal	%d\n", m1.Field)
	fmt.Printf("reflect	%d\n", m2.Field)
}
