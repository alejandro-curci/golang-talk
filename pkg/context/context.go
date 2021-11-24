package context

import (
	"context"
	"fmt"
	"time"
)

// two purposes:
// 1- to provide an API for cancelling branches of your call-graph (prevent our system to do unnecessary work)
// 2- to provide a data-bag for transporting request-scoped data through your call-graph

func summary() {
	// create an empty context
	ctx := context.Background()
	_ = context.TODO()

	// create a child context from a parent context
	ctx, cancel := context.WithCancel(ctx)
	ctx, cancel = context.WithTimeout(ctx, 300*time.Millisecond)
	ctx, cancel = context.WithDeadline(ctx, time.Date(2022, 04, 18, 0, 0, 0, 0, time.UTC))
	cancel() // cancels the new context and all its children

	// context cancellation -> two sides:
	// listening to the event -> wait on "<-ctx.Done()"
	// emitting the signal -> call "cancel()" or set a timeout

	// saving data in the context
	key := "request_id"
	value := 71917101
	ctx = context.WithValue(ctx, key, value)
	// reading the value (cannot be modified, a new context needs to be created instead)
	fmt.Printf("ID: %v", ctx.Value(key))

}

// NODE-TREE STRUCTURE
// Parent and child contexts build a node tree structure, cancelling a parent cancels every child as well.
// However, there’s nothing that allows the function accepting the context to cancel it.
// This protects functions up the call stack from children canceling the context.

// DON'T USE CONTEXT INSTANCES AS STRUCT FIELDS
// Instances of context.Context may look equivalent from the out-side,
// but internally they may change at every stack-frame.
// For this reason, it’s important to always pass instances of context into your functions.

// DON'T USE CONTEXT VALUES FOR PASSING OPTIONAL PARAMETERS
// Use context values only for request-scoped data, not for passing optional parameters to functions.
// suggestions -> request id, user id, authorization token
