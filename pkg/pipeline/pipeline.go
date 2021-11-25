package pipeline

import "fmt"

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
	for n := range power(generate(2, 9, 17, 20, 5, 31)) {
		fmt.Println(n)
	}

	fmt.Println(<-sum(power(power(generate(1, 7, 92, 45)))))
}
