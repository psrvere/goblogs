package memorymodel

import (
	"fmt"
	"sync"
	"time"
)

// Link - https://go.dev/ref/mem

// The go statement that starts a new goroutine is synchronized before the start of the
// goroutine's execution
// In following example go routine will print "bye world" because the main goroutine will
// be synchronized (executed fully) before execution of the child goroutine
func Example1() {
	var a string = "hello world"
	go func() {
		fmt.Println(a)
	}()
	a = "bye world"
	time.Sleep(1 * time.Second) // wait to let goroutine print value of a
}

// The exit of a goroutine is NOT guranteed to be synchronized before any event in the program
// Here a will be printed as "empty" because assignment by go routine is not followed by any
// synchronization event
// In this example, an aggressive compiler might just delete the whole go statement
func Example2() {
	var a string = "empty"
	go func() { a = "hello world" }()
	fmt.Println(a)
}

// A. Channel Communication
// Channel communication is the main method of synchronization between goroutines

// A send on a buffered channel is synchronized before the completion of the corresponding receive
// from than channel
// Here a is guaranteed to print as "hello world"
func Example3() {
	c := make(chan int, 10) // buffered channel
	var a string

	go func() {
		a = "hello world"
		c <- 0
	}()
	<-c
	fmt.Println(a)
}

// Similar to Example3, the closing of a channel is synchronized before a receive
// close(c) and c <- 0 are equivalent statements in these examples where c is a buffered channel
func Example4() {
	c := make(chan int, 10) // buffered channel
	var a string

	go func() {
		a = "hello world"
		close(c)
	}()
	fmt.Println(<-c) // closed channel send 0 value
	fmt.Println(<-c) // recurring calls on a closed channel returns 0
	fmt.Println(a)
}

// A receive from an unbuffered channel is synchronized before the completion of the
// corresponding send on that channel
func Example5() {
	var c = make(chan int) // unbuffered channel
	var a string
	go func() {
		a = "hello world"
		<-c
	}()
	c <- 0
	fmt.Println(a)
}

type FuncType func(x int)

// The kth receive on a channel with capacity C is synchronized before the completion of the
// k+Cth send from that channel completes.
func Example6() {
	work := make([]FuncType, 10)

	for i := range work {
		work[i] = func(x int) { fmt.Println("work function: ", x) }
	}

	limit := make(chan int, 3) // buffered channel
	for i, w := range work {
		go func(w FuncType) {
			// goroutines coordinate using the limit channel to ensure at any given point
			// there are at most 3 work functions running
			limit <- 1
			w(i)
			<-limit
		}(w)
	}
	time.Sleep(1 * time.Second)
}

// B. Locks
// sync.Mutex and sync.RWMutex

// For any sync.Mutex or sync.RWMutex variable l and n < m
// call n o l.Unlock() is synchronized before call m of l.Lock() returns
func Example7() {
	var l sync.Mutex
	var a string

	var f = func() {
		a = "hello world"
		l.Unlock() // the first call to l.Unlock() is synchronized before the second call to l.Lock() returns
	}

	l.Lock()
	go f()
	l.Lock()       // the second call to l.Lock() is sequenced before the print statement
	fmt.Println(a) // a is guaranteed to print "hello world"
}
