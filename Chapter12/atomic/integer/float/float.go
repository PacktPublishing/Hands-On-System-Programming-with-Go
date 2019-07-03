package main

import (
	"fmt"
	"math"
	"sync"
	"sync/atomic"
)

type f64 uint64

func uf(u uint64) (f float64) { return math.Float64frombits(u) }
func fu(f float64) (u uint64) { return math.Float64bits(f) }

func (f *f64) Load() float64 {
	return uf(atomic.LoadUint64((*uint64)(f)))
}

func (f *f64) Store(s float64) {
	atomic.StoreUint64((*uint64)(f), fu(s))
}

func newF64(f float64) *f64 {
	v := f64(fu(f))
	return &v
}

func (f *f64) Add(s float64) float64 {
	for {
		old := f.Load()
		new := old + s
		if f.CompareAndSwap(old, new) {
			return new
		}
	}
}

func (f *f64) CompareAndSwap(old, new float64) bool {
	return atomic.CompareAndSwapUint64((*uint64)(f), fu(old), fu(new))
}

func main() {
	f := newF64(0.54)
	wg := sync.WaitGroup{}
	// 2*iteration + reset at 5
	wg.Add(20)
	for i := 0; i < 10; i++ {
		go func() {
			f.Add(0.01)
			fmt.Println("add")
			wg.Done()
		}()
		go func() {
			fmt.Println("load", f.Load())
			wg.Done()
		}()
	}
	wg.Wait()
}
