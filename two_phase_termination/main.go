// Copyright 2012 T.Yokoyama All rights reserved.

package main

import (
	"fmt"
	"time"
)

type CountUpThread struct {
	counter int64
	shutdownRequested bool
}

func(thread *CountUpThread) ShutdownRequest() {
	thread.shutdownRequested = true
	// interrupt()
}

func(thread CountUpThread) IsShutdownRequested() bool {
	return thread.shutdownRequested
}

func(thread *CountUpThread) Start(ch chan int) {
	for !thread.IsShutdownRequested() {
		thread.doWork()
	}

	thread.doShutdown()
	ch <- 0
}

func(thread *CountUpThread) doWork() {
	thread.counter++
	fmt.Printf("doWork: counter = %d\n", thread.counter)
	time.Sleep(500 * time.Millisecond)
}

func(thread CountUpThread) doShutdown() {
	fmt.Printf("doShutdown: counter = %d\n", thread.counter)
}

func main() {
	fmt.Println("main: BEGIN")

	ch := make(chan int)
	thread := CountUpThread{counter: 0, shutdownRequested: false}
	go thread.Start(ch)

	time.Sleep(10 * time.Second)

	fmt.Println("main: shutdownRequest")
	thread.ShutdownRequest()

	// wait CountUpThread is finished.
	fmt.Println("main: join")
	<-ch

	fmt.Println("main: END")
}