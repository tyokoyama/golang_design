// Copyright 2012 T.Yokoyama All rights reserved.

package main

import (
	"container/list"
	"fmt"
	"math/rand"
	"sync"
	"time"
)

type Request struct {
	Name string
}

func(req Request) String() string {
	return fmt.Sprintf("[ Request %s ]", req.Name)
}

type RequestQueue struct {
	queue *list.List
	cond *sync.Cond
}

// 構造体のメソッドはinit()でも呼ばれない！
func (queue *RequestQueue) init() {
	queue.queue = list.New()
	var mutex sync.Mutex
	queue.cond = sync.NewCond(&mutex)
}

func (queue RequestQueue) getRequest() Request {
	queue.cond.L.Lock()
	if queue.queue.Len() <= 0 {
		// wait
		println("getRequest wait")
		queue.cond.Wait()
	}

	element := queue.queue.Front()
	request := element.Value.(Request)
	queue.queue.Remove(element)

	queue.cond.L.Unlock()
	return request
}

func (queue RequestQueue) putRequest(req Request) {
	queue.cond.L.Lock()
	queue.queue.PushBack(req)
	queue.cond.Signal()
	queue.cond.L.Unlock()
}

type ClientThread struct {
	random *rand.Rand
	requestQueue *RequestQueue
}

func (thread ClientThread) Start(result chan int) {
	for i := 0; i < 10000; i++ {
		request := Request {Name: fmt.Sprintf("No. %d", i)}
		fmt.Printf("ClientThread requests %s\n", request)
		thread.requestQueue.putRequest(request)
		time.Sleep(time.Duration(thread.random.Int31()))
	}

	result <- 1
}

type ServerThread struct {
	random *rand.Rand
	requestQueue *RequestQueue
}

func (thread ServerThread) Start(result chan int) {
	for i := 0; i < 10000; i++ {
		request := thread.requestQueue.getRequest()
		fmt.Printf("ServerThread handles %s\n", request)
		time.Sleep(time.Duration(thread.random.Int31()))
	}
	result <- 1
}

func main() {
	var requestQueue RequestQueue
	requestQueue.init()

	var channels [2] chan int
	for i := 0; i < len(channels); i++ {
		channels[i] = make(chan int)
	}

	go ClientThread {random: rand.New(rand.NewSource(time.Now().Unix())), requestQueue: &requestQueue}.Start(channels[0])
	go ServerThread {random: rand.New(rand.NewSource(time.Now().Unix())), requestQueue: &requestQueue}.Start(channels[1])

	for i := 0; i < len(channels); i++ {
		<-channels[i]
	}
}