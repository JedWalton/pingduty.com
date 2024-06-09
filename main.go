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

// func frontend() {
// 	// Serve static assets like CSS, JS, and images
// 	fs := http.FileServer(http.Dir("./frontend")) // Serves the current directory
//
// 	// Handle all requests by serving a file of the same name
// 	http.Handle("/", fs)
//
// 	// Log that the server is starting
// 	log.Println("Listening on :8080...")
//
// 	// Start the server on port 8080
// 	err := http.ListenAndServe(":8080", nil)
// 	if err != nil {
// 		log.Fatal("ListenAndServe: ", err)
// 	}
// }
