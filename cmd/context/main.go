package main

import (
	"context"
	"fmt"
	"time"
)

// waits for 2 seconds and then prints the message
// the process can be cancelled at any time by writing to the console
func main() {
	ctx := context.Background()
	ctx, cancel := context.WithCancel(ctx)

	go func() {
		var s string
		_, _ = fmt.Scanf("%s", &s)
		cancel()
	}()

	mySleepAndTalk(ctx, 5*time.Second, "PRINTING STUFF HERE...")
}

func mySleepAndTalk(ctx context.Context, d time.Duration, msg string) {
	select {
	case <-ctx.Done():
		fmt.Println("ERROR: ", ctx.Err())
	case <-time.After(d):
		fmt.Println(msg)
	}
}
