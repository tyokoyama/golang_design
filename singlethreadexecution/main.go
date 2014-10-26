// Copyright 2012 T.Yokoyama All rights reserved.

package main

import (
	"fmt"
	"math/rand"
	"time"
	"runtime"
)

type Gate struct {
	counter int
	name string
	address string
	ch chan int
}

func(g *Gate) String() string {
	return fmt.Sprintf("No.%d: %s, %s", g.counter, g.name, g.address)
}

func(g *Gate) check() {
//	fmt.Printf("%s\n", g)
	if(g.name[0] != g.address[0]) {
		println("***** BROKEN ***** " + g.String())
	}
}

func(g *Gate) Pass(name, address string) {
	g.ch <- 0
	g.counter++
	g.name = name
	g.address = address
	g.check()
	<- g.ch

}

type UserThread struct {
	gate *Gate
	myName string
	myAddress string
}

func(thread UserThread) Start() {
	println(fmt.Sprintf("%s BEGIN", thread.myName))
	for {
		thread.gate.Pass(thread.myName, thread.myAddress)
		r := (rand.Int() % 10) + 1
//		time.Sleep(1 * time.Millisecond)
		d, _ := time.ParseDuration(fmt.Sprintf("%dms", r))
		time.Sleep(d)
	}
}

func main() {
	println("Testing Gate, hit CTRL+C to exit.")
	ch := make(chan int, 1)
	runtime.GOMAXPROCS(3)

	gate := Gate {counter: 0, name: "Nobody", address: "Nowhere", ch: ch}
	alice := UserThread {gate: &gate, myName: "Alice", myAddress: "Alaska"}
	bobby := UserThread {gate: &gate, myName: "Bobby", myAddress: "Brazil"}
	chris := UserThread {gate: &gate, myName: "Chris", myAddress: "Canada"}

	go alice.Start()
	go bobby.Start()
	go chris.Start()

	select {}
}
