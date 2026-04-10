package main

import (
	"fmt"
	"net/http"
	"os"
)

func main() {
	// Static file serving (for our UI)
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "index.html")
	})

	// The "Engine" endpoint: This is where we trigger your Go functions
	http.HandleFunc("/run", handleCalculation)

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	fmt.Printf("Pi Suite starting on http://localhost:%s\n", port)
	http.ListenAndServe(":"+port, nil)
}

// This replaces your TrafficManager logic
func handleCalculation(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/event-stream")
	w.Header().Set("Cache-Control", "no-cache")
	w.Header().Set("Connection", "keep-alive")

	method := r.URL.Query().Get("method")

	// Create a channel to catch the output strings from your math funcs
	outputChan := make(chan string)
	done := make(chan bool)
	go func() {
		switch method {
		case "archimedes":
			// Only one argument here now
			ArchimedesBig()
		case "roots":
			runRootsWeb(r.URL.Query(), func(s string) { outputChan <- s })
		}
		// This line tells the web browser the calculation is finished
		done <- true
	}()

	// Stream the data to the browser
	for {
		select {
		case msg := <-outputChan:
			fmt.Fprintf(w, "data: %s\n\n", msg)
			w.(http.Flusher).Flush()
		case <-done:
			return
		case <-r.Context().Done():
			return
		}
	}
}
