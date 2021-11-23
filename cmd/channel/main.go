package main

import (
	"fmt"
)

func main() {
	example1()
}

func right_way() {
	stream := make(chan int)
	go func() {
		defer close(stream)
		for i := 0; i < 5; i++ {
			stream <- i
		}
	}()
	for n := range stream {
		fmt.Println(n)
	}
}

// deadlock!
// the range loop can't read anymore because the goroutine is not sending anything
func example1() {
	stream := make(chan int)
	go func() {
		for i := 0; i < 5; i++ {
			stream <- i
		}
	}()
	for n := range stream {
		fmt.Println(n)
	}
}

// no deadlock, but prints infinite "0 false"
func example2() {
	stream := make(chan int)
	go func() {
		defer close(stream)
		for i := 0; i < 5; i++ {
			stream <- i
		}
	}()
	// keeps reading forever!
	for {
		n, ok := <-stream
		fmt.Printf("%v %v\n", n, ok)
	}
}
