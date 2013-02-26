// Copyright 2012 T.Yokoyama All rights reserved.

package main

import (
	"fmt"
	"sync"
	"time"
)

type Data interface {
	GetContent() string
}

type FutureData struct {
	realData RealData
	ready bool
	cond *sync.Cond
}

func(data *FutureData) init() {
	var mutex sync.Mutex
	data.cond = sync.NewCond(&mutex)
	data.ready = false
}

func(data *FutureData) SetRealData(realData RealData) {
	data.cond.L.Lock()
	defer data.cond.L.Unlock()
	if(data.ready) {
		return		// balk
	}

	data.realData = realData
	data.ready = true
	data.cond.Signal()
}

func(data *FutureData) GetContent() string {
	data.cond.L.Lock()
	for !data.ready {
		data.cond.Wait()
	}

	content := data.realData.GetContent()
	data.cond.L.Unlock()

	return content
}

type RealData struct {
	content string
}

func(data *RealData) Make(count int, c byte) {
	fmt.Printf("        making RealData(%d, %c) BEGIN\n", count, c)
	buffer := make([]byte, count)
	for i := 0; i < count; i++ {
		buffer[i] = c
		time.Sleep(100 * time.Millisecond)
	}
	fmt.Printf("        making RealData(%d, %c) END\n", count, c)
	data.content = string(buffer)
}

func(data *RealData) GetContent() string {
	return data.content
}

type Host struct {

}

func(host Host) request(count int, c byte) Data {
	fmt.Printf("    request(%d, %c) BEGIN\n", count, c)

	var future FutureData
	future.init()
	go func() {
		var realData RealData
		realData.Make(count, c)
		future.SetRealData(realData)
	}()

	fmt.Printf("    request(%d, %c) END\n", count, c)

	return &future
}

func main() {
	fmt.Printf("main BEGIN\n")

	var host Host
	data1 := host.request(10, 'A')
	data2 := host.request(20, 'B')
	data3 := host.request(30, 'C')

	fmt.Printf("main otherJob BEGIN\n")
	time.Sleep(2 * time.Second)
	fmt.Printf("main otherJob END\n")

	fmt.Printf("data1 = %s\n", data1.GetContent())
	fmt.Printf("data2 = %s\n", data2.GetContent())
	fmt.Printf("data3 = %s\n", data3.GetContent())

	fmt.Printf("main END\n")
}