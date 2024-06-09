package frontend

import (
	"log"
	"net/http"
)

func Run() {
	// Serve static assets like CSS, JS, and images
	fs := http.FileServer(http.Dir("./frontend")) // Serves the current directory

	// Handle all requests by serving a file of the same name
	http.Handle("/", fs)

	// Log that the server is starting
	log.Println("Listening on :8080...")

	// Start the server on port 8080
	err := http.ListenAndServe(":8080", nil)
	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
}
