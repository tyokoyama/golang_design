// Copyright 2012 T.Yokoyama All rights reserved.

package main

import (
	"fmt"
	"time"
)

type Person struct {
	Name string
	Address string
}

func(p Person) String() string {
	return fmt.Sprintf("[ Person: name = %s, address = %s ]", p.Name, p.Address)
}

type PrintPersonThread struct {
	person Person
	id int
}

func(t PrintPersonThread) Start() {
	for {
		println(fmt.Sprintf("%d prints %s", t.id, t.person))
		time.Sleep(1 * time.Millisecond)
	}
}

func main() {
	alice := Person {Name: "Alice", Address: "Alaska"}

	go PrintPersonThread {id: 1, person: alice}.Start()
	go PrintPersonThread {id: 2, person: alice}.Start()
	go PrintPersonThread {id: 3, person: alice}.Start()

	select {}
}