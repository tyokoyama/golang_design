// Copyright 2012 T.Yokoyama All rights reserved.

package main

import (
	"fmt"
	"math/rand"
	"os"
	"sync"
	"runtime"
	"time"
)

type Data struct {
	filename string	// ファイル名
	content string	// データの内容
	changed bool	// 変更の有無
	mutex sync.Mutex
	// mutex sync.RWMutex
}

func (data *Data) Change(newContent string) {
	data.mutex.Lock()
	defer data.mutex.Unlock()
	data.content = newContent
	data.changed = true
}

func (data *Data) Save(name string) {
	// Golang Cafe #49: 書き込み処理をgoroutine化して、更新フラグを参照するのはガンガンいけばいいのでは？
	data.mutex.Lock()
	defer data.mutex.Unlock()

	if(!data.changed) {
		return;
	}

	data.doSave(name)
	data.changed = false
}

func (data Data)doSave(name string) {
	fmt.Printf("%s Called doSave, content = %s\n", name, data.content)
	file, err := os.OpenFile(data.filename, os.O_WRONLY|os.O_CREATE, os.ModePerm)
	if err != nil {
		panic(err.Error())
	}
	defer file.Close()

	if _, err := fmt.Fprintf(file, "%s", data.content); err != nil {
		panic(err.Error())
	}

	// time.Sleep(10 * time.Second)
}

type SaverThread struct {
	data *Data
}

func (thread SaverThread) Start() {
	for {
		thread.data.Save("SaverThread")
		time.Sleep(1 * time.Second)
	}
}

type ChangerThread struct {
	data *Data
	random *rand.Rand
}

func (thread ChangerThread) Start() {
	for i := 0; ; i++ {
		thread.data.Change(fmt.Sprintf("No. %d\n", i))
		time.Sleep(time.Duration(thread.random.Int31()))
		thread.data.Save("ChangerThread")
	}
}

func main() {
	runtime.GOMAXPROCS(runtime.NumCPU())

	data := &Data {filename: "data.txt", content: "(empty)", changed: true}
	go ChangerThread {data: data, random: rand.New(rand.NewSource(time.Now().Unix()))}.Start()
	go SaverThread {data: data}.Start()

	select {}
}