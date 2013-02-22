// Copyright 2012 T.Yokoyama All rights reserved.

package main

import (
	"fmt"
	"math/rand"
	"sync"
	"time"
)

type Table struct {
	buffer []string
	tail int
	head int
	count int
	cond *sync.Cond
}

func (t *Table) init() {
	var mutex sync.Mutex
	t.cond = sync.NewCond(&mutex)
}

func (t *Table) put(cake, name string) {
	t.cond.L.Lock()
	fmt.Printf("%s puts %s\n", name, cake)
	for t.count >= len(t.buffer) {
		// wait
		t.cond.Wait()
	}

	t.buffer[t.tail] = cake
	t.tail = (t.tail + 1) % len(t.buffer)
	t.count++
	t.cond.Signal()
	t.cond.L.Unlock()
}

func (t *Table) take(name string) string {
	t.cond.L.Lock()
	for t.count <= 0 {
		// wait
		t.cond.Wait()
	}

	cake := t.buffer[t.head]
	t.head = (t.head + 1) % len(t.buffer)
	t.count--

	fmt.Printf("%s takes %s\n", name, cake)
	t.cond.Signal()
	t.cond.L.Unlock()
	return cake
}

var id int = 0
var id_mutex sync.Mutex

func nextId() int {
	id_mutex.Lock()
	id++
	nextId := id
	id_mutex.Unlock()
	return nextId
}

type MakerThread struct {
	random *rand.Rand
	table *Table
	name string
}

func (thread MakerThread) Start() {
	for {
		time.Sleep(time.Duration(thread.random.Int31()))
		cake := fmt.Sprintf("[ Cake No. %d by %s ]", nextId(), thread.name)
		thread.table.put(cake, thread.name)
	}
}

type EaterThread struct {
	random *rand.Rand
	table *Table	
	name string
}

func (thread EaterThread) Start() {
	for {
		thread.table.take(thread.name)
		time.Sleep(time.Duration(thread.random.Int31()))
	}
}

func main() {
	buffer := make([]string, 3)
	table := Table {buffer: buffer, head: 0, tail: 0, count: 0}
	table.init()

	go MakerThread {random: rand.New(rand.NewSource(time.Now().Unix())), table: &table, name: "MakerThread-1"}.Start()
	go MakerThread {random: rand.New(rand.NewSource(time.Now().Unix())), table: &table, name: "MakerThread-2"}.Start()
	go MakerThread {random: rand.New(rand.NewSource(time.Now().Unix())), table: &table, name: "MakerThread-3"}.Start()

	go EaterThread {random: rand.New(rand.NewSource(time.Now().Unix())), table: &table, name: "EaterThread-1"}.Start()
	go EaterThread {random: rand.New(rand.NewSource(time.Now().Unix())), table: &table, name: "EaterThread-2"}.Start()
	go EaterThread {random: rand.New(rand.NewSource(time.Now().Unix())), table: &table, name: "EaterThread-3"}.Start()

	select {}
}
