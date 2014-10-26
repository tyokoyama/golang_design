package main

import (
	"fmt"
	"runtime"
	"time"
)

type Host struct {
	ch []chan int
	Helper Helper
}

func(h *Host) Request(count int, c byte) {
	ch := make(chan int)
	fmt.Printf("    request(%d, %c) BEGIN\n", count, c)

	go func(ch chan<- int) {
		h.Helper.Handle(count, c)

		ch <- count
	}(ch)

	fmt.Printf("    request(%d, %c) END\n", count, c)

	h.ch = append(h.ch, ch)
}

func (h Host) Wait() {
	for i := 0; i < len(h.ch); i++ {
		// fmt.Println("Wait", <-h.ch[i])
		<-h.ch[i]
	}
}

type Helper bool

func(h Helper) Handle(count int, c byte) {
	fmt.Printf("        handle(%d, %c) BEGIN\n", count, c)

	for i := 0; i < count; i++ {
		time.Sleep(100 * time.Millisecond)
		fmt.Printf("%c", c)
	}

	fmt.Println("")
	fmt.Printf("        handle(%d, %c) END\n", count, c)
}

func main() {
	var h Host

	runtime.GOMAXPROCS(runtime.NumCPU())
	fmt.Println("main BEGIN")

	h.Request(10, 'A')
	h.Request(20, 'B')
	h.Request(30, 'C')

	fmt.Println("main END")

	h.Wait()
}