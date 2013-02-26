// Copyright 2012 T.Yokoyama All rights reserved.

package main

import (
	"fmt"
	"math/rand"
	"sync"
	"time"
)

type Request struct {
	name string
	number int
	random *rand.Rand
}

func (req Request) execute(name string) {
	fmt.Printf("%s executes %s\n", name, req)
}

func (req Request) String() string {
	return fmt.Sprintf("[ Request from %s No. %d ]", req.name, req.number)
}

type WorkerThread struct {
	channel *Channel
}
func (thread WorkerThread) Start(name string) {
	for {
		request := thread.channel.takeRequest()
		request.execute(name)
	}
}

const MAX_REQUEST = 100
type Channel struct {
	requestQueue []Request
	tail int
	head int
	count int
	threadPool []WorkerThread
	cond *sync.Cond
}

func(c *Channel) init(threads int) {
	c.requestQueue = make([]Request, MAX_REQUEST)
	c.threadPool = make([]WorkerThread, threads)
	for i := 0; i < len(c.threadPool); i++ {
		c.threadPool[i].channel = c
	}
	var mutex sync.Mutex
	c.cond = sync.NewCond(&mutex)
}

func(c *Channel) startWorkers() {
	for i := 0; i < len(c.threadPool); i++ {
		go c.threadPool[i].Start(fmt.Sprintf("Worker-%d", i))
	}
}

func(c *Channel) putRequest(request Request) {
	c.cond.L.Lock()
	for c.count >= len(c.requestQueue) {
		c.cond.Wait()
	}

	c.requestQueue[c.tail] = request
	c.tail = (c.tail + 1) % len(c.requestQueue)
	c.count++
	c.cond.Signal()
	c.cond.L.Unlock()
}

func(c *Channel) takeRequest() Request {
	c.cond.L.Lock()

	for c.count <= 0 {
		c.cond.Wait()
	}

	request := c.requestQueue[c.head]
	c.head = (c.head + 1) % len(c.requestQueue)
	c.count--

	c.cond.Signal()
	c.cond.L.Unlock()

	return request
}

type ClientThread struct {
	channel *Channel
	random *rand.Rand
}

func (thread ClientThread) Start(name string) {
	for i := 0; ; i++ {
		request := Request {name: name, number: i}
		thread.channel.putRequest(request)
		time.Sleep(time.Duration(thread.random.Int31()))
	}
}

func main() {
	channel := Channel{tail: 0, head: 0, count: 0}
	channel.init(5)
	channel.startWorkers()
	go ClientThread{channel: &channel, random: rand.New(rand.NewSource(time.Now().Unix()))}.Start("Alice")
	go ClientThread{channel: &channel, random: rand.New(rand.NewSource(time.Now().Unix()))}.Start("Bobby")
	go ClientThread{channel: &channel, random: rand.New(rand.NewSource(time.Now().Unix()))}.Start("Chris")

	select {}
}