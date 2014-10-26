package main

import (
	"fmt"
	"runtime"
	"sync"
	"time"
)

func main() {
	runtime.GOMAXPROCS(runtime.NumCPU())

	var p sync.Pool

	p.New = func() interface{} {
		fmt.Println("p.New execute.")
		time.Sleep(100 * time.Millisecond)

		return true
	}

	p.Put(1)
	p.Put(2)
	p.Put(3)

	fmt.Println(p.Get())
	fmt.Println("1 Get.")
	fmt.Println(p.Get())
	fmt.Println("2 Get.")
	fmt.Println(p.Get())
	fmt.Println("3 Get.")
}
