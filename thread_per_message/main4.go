package main

import (
	"fmt"
//	"runtime"
	"sync"
)

var p sync.Pool

type Param struct {
	Value int
	IsClose bool
	IsPut bool
}

func main() {
//	runtime.GOMAXPROCS(runtime.NumCPU())

	// Get()した時にPutしたデータが1件もなかった時に呼び出される。
	// 定義していなければnilが返る。
	p.New = func () interface{} {
		return interface{}(-1)
	}

	st := []Param {
		Param {Value:1, IsClose: false, IsPut: true},
		Param {Value:2, IsClose: false, IsPut: true},
		Param {Value:3, IsClose: false, IsPut: false},
		Param {Value:4, IsClose: true, IsPut: true},
	}

	for _, s := range st {
		if s.IsPut {
			p.Put(s.Value)
		}
	}

	ch := make(chan int, 3)

	for pos, s := range st {
		go printValue(pos + 1, ch, s.IsClose)
	}

	for res := range ch {
		fmt.Println("main: ", res)
	}
}

func printValue(i int, ch chan int, isClose bool) {

	pv := p.Get()

	fmt.Printf("printValue[%d] = %d\n", i, pv)

	ch <- pv.(int)

	if isClose {
		close(ch)
	}
}