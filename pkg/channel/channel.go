package channel

import (
	"fmt"
	"time"
)

// "Do not communicate by sharing memory; instead, share memory by communicating." (Rob Pike)

func summary() {
	// channels communicate information between goroutines
	// they can act either as emitters or receivers
	ch := make(chan int)  // unbuffered (no capacity)
	_ = make(chan int, 3) // buffered with capacity of 3

	// channels are blocking
	// a send to a full channel will wait until it is blocked
	// a read from an empty channel will wait until one item is placed in it
	go func() {
		ch <- 17 // send to a channel
	}()
	fmt.Println(<-ch) // read from a channel

	// this is wrong! it causes a deadlock
	ch <- 3
	fmt.Printf("received %v", <-ch)

	// channels can be closed to indicate no more values will be sent over it
	// send to a closed channel will cause a panic
	// read from a closed channel will produce the zero-value from the channel type
	close(ch)
	num, ok := <-ch // 0 false
	fmt.Printf("Zero value: %v and boolean: %v", num, ok)

	// ranging over a channel and closing to break the loop
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

	// the select statement allows us to handle more than one channel input within a goroutine
	// all channels are considered simultaneously
	// if none of the channels are ready, the select statement blocks
	// when one channel is ready, that operation proceeds
	c := make(chan string)
	go func() {
		// make a request
		time.Sleep(1*time.Second)
		c <- "response"
	}()
	select {
	case <-c:
		fmt.Println("request processed")
	case <-time.After(500 * time.Millisecond):
		fmt.Println("timeout!")
	}

}
