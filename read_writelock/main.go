// Copyright 2012 T.Yokoyama All rights reserved.

package main

import (
	"fmt"
	"math/rand"
	"sync"
	"time"
)

type ReadWriteLock struct {
	readingReaders int
	waitingWriters int
	writingWriters int
	preferWriter bool

	cond *sync.Cond
}

func (lock *ReadWriteLock) init() {
	var mutex sync.Mutex
	lock.cond = sync.NewCond(&mutex)
}

func (lock *ReadWriteLock) readLock() {
	lock.cond.L.Lock()
	for lock.writingWriters > 0 || (lock.preferWriter && lock.waitingWriters > 0) {
		lock.cond.Wait()
	}

	lock.readingReaders++
}

func (lock *ReadWriteLock) readUnlock() {
	lock.readingReaders--
	lock.preferWriter = true
	lock.cond.Signal()
	lock.cond.L.Unlock()
}

func (lock *ReadWriteLock) writeLock() {
	lock.cond.L.Lock()	
	lock.waitingWriters++
	for lock.readingReaders > 0 || lock.writingWriters > 0 {
		lock.cond.Wait()
	}

	lock.waitingWriters--

	lock.writingWriters++
}

func (lock *ReadWriteLock) writeUnLock() {
	lock.writingWriters--
	lock.preferWriter = false
	lock.cond.Signal()
	lock.cond.L.Unlock()
}

type Data struct {
	buffer string
	size int
	lock ReadWriteLock
}

func(data *Data) init(size int) {
	for i := 0; i < size; i++ {
		data.buffer += "*"
	}
	data.size = size
}

func(data *Data) Read() string {
	data.lock.readLock()
	result := data.doRead()
	data.lock.readUnlock()

	return result
}

func(data *Data) Write(c byte) {
	data.lock.writeLock()
	data.doWrite(c)
	data.lock.writeUnLock()
}

func(data *Data) doRead() string {
	result := data.buffer
	data.slowly()
	return result
}

func(data *Data) doWrite(c byte) {
	data.buffer = ""
	for i := 0; i < data.size; i++ {
		data.buffer += string(c)
		data.slowly()
	}
}

func(data *Data) slowly() {
	time.Sleep(50 * time.Millisecond)
}

type WriterThread struct {
	random *rand.Rand
	data *Data
	filler string
	index int
}

func(thread WriterThread) Start() {
	for {
		str := thread.nextchar()
		thread.data.Write(str)
		time.Sleep(time.Duration(thread.random.Int31()))
	}
}

func(thread *WriterThread) nextchar() byte {
	str := thread.filler[thread.index]
	thread.index++
	if thread.index >= len(thread.filler) {
		thread.index = 0
	}
	return str
}

type ReaderThread struct {
	data *Data
}

func(thread ReaderThread) Start(name string) {
	for {
		str := thread.data.Read()
		fmt.Printf("%s reads %s\n", name, str)
	}
}

func main() {
	lock := ReadWriteLock {readingReaders: 0, writingWriters: 0, waitingWriters: 0, preferWriter: true}
	lock.init()
	data := Data {lock: lock}
	data.init(10)

	go ReaderThread {data: &data}.Start("Thread-0")
	go ReaderThread {data: &data}.Start("Thread-1")
	go ReaderThread {data: &data}.Start("Thread-2")
	go ReaderThread {data: &data}.Start("Thread-3")
	go ReaderThread {data: &data}.Start("Thread-4")
	go ReaderThread {data: &data}.Start("Thread-5")
	go WriterThread {random: rand.New(rand.NewSource(time.Now().Unix())), data: &data, filler: "ABCDEFGHIJKLMNOPQRSTUVWXYZ", index: 0}.Start()
	go WriterThread {random: rand.New(rand.NewSource(time.Now().Unix())), data: &data, filler: "abcdefghijklmnopqrstuvwxyz", index: 0}.Start()	

	select {}
}