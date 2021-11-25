package pipeline

import (
	"fmt"
	"sync"
)

// A pipeline is a series of stages connected by channels,
// where each stage is a group of goroutines running the same function.

// In each stage, the goroutines:
// 1- receive values from upstream via inbound channels
// 2- perform some function on that data, producing new values
// 3- send those values downstream via outbound channels

// first stage -> source or producer
// last stage -> sink or consumer

// generate is the first stage, it converts a list of integers into a channel which emits those numbers
func generate(numbers ...int) <-chan int {
	out := make(chan int)
	go func() {
		for _, n := range numbers {
			out <- n
		}
		close(out)
	}()
	return out
}

// power is the second stage, it powers the numbers received from stage 1 and sends them to another channel
func power(in <-chan int) <-chan int {
	out := make(chan int)
	go func() {
		for n := range in {
			out <- n * n
		}
		close(out)
	}()
	return out
}

// sum is the third stage, it sums all the numbers received from the previous stage
func sum(in <-chan int) <-chan int {
	out := make(chan int)
	go func() {
		var total int
		for n := range in {
			total += n
		}
		out <- total
		close(out)
	}()
	return out
}

// main client, it consumes the numbers coming from a pipeline (printing to console)
func main() {

	// before
	c := generate(2, 9, 17, 20, 5, 31)
	for n := range power(c) {
		fmt.Println(n)
	}

	// fan-out/fan-in
	in := generate(2, 9, 17, 20, 5, 31)
	// distribute the power() work across three goroutines (reading from the same channel in)
	ch1 := power(in)
	ch2 := power(in)
	ch3 := power(in)

	for n := range merge(ch1, ch2, ch3) {
		fmt.Println(n)
	}
}

// FAN-OUT
// Multiple functions can read from the same channel until that channel is closed.
// It provides a way to distribute work amongst a group of workers to parallelize CPU use and I/O.

// FAN-IN
// A function can read from multiple inputs and proceed until all are closed by multiplexing
// the input channels onto a single channel that’s closed when all the inputs are closed.

func merge(channels ...<-chan int) <-chan int {
	var wg sync.WaitGroup
	out := make(chan int)

	// closure -> sends values from channels into the out channel
	send := func(ch <-chan int) {
		for n := range ch {
			out <- n
		}
		wg.Done()
	}

	wg.Add(len(channels))
	for _, ch := range channels {
		go send(ch)
	}

	// wait and close the out channel in a different goroutine
	go func() {
		wg.Wait()
		close(out)
	}()

	return out
}
