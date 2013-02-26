// Copyright 2012 T.Yokoyama All rights reserved.

package main

import (
	"fmt"
	"time"
)

type Helper struct {

}

func(h Helper) Handle(count int, c byte) {
	fmt.Printf("        handle(%d, %c) BEGIN\n", count, c)
	for i := 0; i < count; i++ {
		h.slowly()
		fmt.Printf("%c", c)
	}
	fmt.Println("")
	fmt.Printf("        handle(%d, %c) END\n", count, c)
}

func(h Helper) slowly() {
	time.Sleep(100 * time.Millisecond)
}

type Host struct {
	Helper Helper
}

func(h Host) request(count int, c byte) chan bool {
	fmt.Printf("    request(%d, %c) BEGIN\n", count, c)
	ch := make(chan bool)
	go func() {
		h.Helper.Handle(count, c)
		ch <- true
	}()
	fmt.Printf("    request(%d, %c) END\n", count, c)

	return ch
}

func main() {
	fmt.Printf("main BEGIN\n")


	var helper Helper
	host := Host{Helper: helper}
	ch1 := host.request(10, 'A')
	ch2 := host.request(20, 'B')
	ch3 := host.request(30, 'C')
	fmt.Printf("main END\n")

	<-ch1
	<-ch2
	<-ch3
}