package main

import (
	"sync"

	"pingduty.com/api/uptime"
	"pingduty.com/frontend"
)

func main() {
	var wg sync.WaitGroup
	// Increment the WaitGroup counter for each concurrent task
	wg.Add(2)

	// Start frontend server in a goroutine
	go func() {
		defer wg.Done() // Decrement the counter when the goroutine completes
		frontend.Run()
	}()

	// Start pingduty function in a separate goroutine
	go func() {
		defer wg.Done()
		uptime.Run()
	}()

	// Wait for all goroutines to complete
	wg.Wait()
}
