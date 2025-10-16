package main

import (
	"sync"
)

// Manufacturing a deadlock by forcing both goroutines to wait until they have
// both acquired a lock...

func main() {
	var a, b sync.Mutex
	var wg sync.WaitGroup
	wg.Add(2)

	ready := make(chan struct{}, 2) // signal: "I acquired my first lock"
	start := make(chan struct{})    // barrier: "now try for the second lock"

	// Goroutine 1: lock A, then (after barrier) try B
	go func() {
		a.Lock()
		ready <- struct{}{} // tell main we've got A
		<-start             // wait until both goroutines hold their first lock
		b.Lock()            // <- deadlock here
		defer b.Unlock()
		defer a.Unlock()
		wg.Done()
	}()

	// Goroutine 2: lock B, then (after barrier) try A
	go func() {
		b.Lock()
		ready <- struct{}{} // tell main we've got B
		<-start
		a.Lock() // <- deadlock here
		defer a.Unlock()
		defer b.Unlock()
		wg.Done()
	}()

	// Wait until both have their first lock
	<-ready
	<-ready
	close(start) // release barrier; both now try to grab the other lock

	wg.Wait() // runtime will panic: "all goroutines are asleep - deadlock!"
}
