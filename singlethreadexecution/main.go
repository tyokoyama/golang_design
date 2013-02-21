// Copyright 2012 T.Yokoyama All rights reserved.

package main

import (
	"fmt"
	"time"
)

type Gate struct {
	counter int
	name string
	address string
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
	g.counter++
	g.name = name
	g.address = address
	g.check()
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
		time.Sleep(1 * time.Millisecond)
	}
}

func main() {
	println("Testing Gate, hit CTRL+C to exit.")

	gate := Gate {counter: 0, name: "Nobody", address: "Nowhere"}
	alice := UserThread {gate: &gate, myName: "Alice", myAddress: "Alaska"}
	bobby := UserThread {gate: &gate, myName: "Bobby", myAddress: "Brazil"}
	chris := UserThread {gate: &gate, myName: "Chris", myAddress: "Canada"}

	go alice.Start()
	go bobby.Start()
	go chris.Start()

	select {}
}
