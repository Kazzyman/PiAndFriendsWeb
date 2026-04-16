package main

import (
	"fmt"
	"net/http"
	"os"
	"strconv"
	"strings"
)

func main() {
	// Static file serving (for our UI)
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "index.html")
	})

	// The "Engine" endpoint: This is where we trigger the Go functions
	http.HandleFunc("/run", handleCalculation)

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	fmt.Printf("Pi Suite starting on http://localhost:%s\n", port)
	http.ListenAndServe(":"+port, nil)
}

// This replaces the TrafficManager logic from the Feyn version.
func handleCalculation(w http.ResponseWriter, r *http.Request) {
	// 1. Set Web Headers
	w.Header().Set("Content-Type", "text/event-stream")
	w.Header().Set("Cache-Control", "no-cache")
	w.Header().Set("Connection", "keep-alive")

	method := r.URL.Query().Get("method")

	// 2. Define the channels (The "Plumbing")
	outputChan := make(chan string)
	done := make(chan bool)

	// 3. Launch the Math Engine in a goroutine
	// This must be inside the go func() where outputChan was defined
	go func() {
		switch method { // ::: These cases are where the association is made to webPrint.
		case "archimedes":
			// Here, we create an 'anonymous' function on the fly. Anything ArchimedesBig sends to 's' gets thrown onto the channel.
			ArchimedesBig(func(s string) {
				outputChan <- s
			})
		case "spigot":
			TheSpigotWeb(done, func(s string) {
				outputChan <- s
			})
		case "monte":
			// We grab the gridSize from the URL, or default to "100" if empty
			gridSize := r.URL.Query().Get("gridSize")
			if gridSize == "" {
				gridSize = "5000"
			}
			MonteCarloWeb(gridSize, func(s string) {
				outputChan <- s
			})
		case "roots":
			runRootsWeb(r.URL.Query(), func(s string) {
				outputChan <- s
			})
		case "bbp":
			digitsStr := r.URL.Query().Get("digits")
			digits, err := strconv.Atoi(digitsStr)
			if err != nil || digits <= 0 {
				digits = 100 // Default fallback
			}

			// Updated to match: func(webPrint, digits, done)
			bbpFast44(func(s string) {
				outputChan <- s
			}, digits, done)
		case "chudnovsky":
			digitsStr := r.URL.Query().Get("digits")
			digits, err := strconv.Atoi(digitsStr)
			if err != nil || digits <= 0 {
				digits = 1000 // Default to a decent number of digits
			}

			// Matches: func(webPrint, digits, done)
			chudnovskyBig(func(s string) {
				outputChan <- s
			}, digits, done)
		case "customseries":
			CustomSeries(done, func(s string) {
				outputChan <- s
			})
		case "wallis":
			JohnWallis(done, func(s string) {
				outputChan <- s
			})
		case "gauss":
			Gauss_Legendre(func(s string) {
				outputChan <- s
			})
			// main.go - inside the switch method block
		case "gregory":
			// Matches: func(webPrint, done)
			GregoryLeibniz(func(s string) {
				outputChan <- s
			}, done)

		case "nilakantha":
			itersStr := r.URL.Query().Get("iters")
			iters, err := strconv.Atoi(itersStr)
			if err != nil || iters <= 0 {
				iters = 1000000
			}
			precStr := r.URL.Query().Get("precision")
			precision, err := strconv.Atoi(precStr)
			if err != nil || precision <= 0 {
				precision = 512
			}
			NilakanthaBig(iters, precision, done, func(s string) {
				outputChan <- s
			})

		case "nilakantha_classic":
			nifty_scoreBoardWeb(done, func(s string) {
				outputChan <- s
			})

		default:
			outputChan <- "Unknown method requested."
		}

		func() {
			defer func() { recover() }()
			done <- true
		}()
	}()

	// 4. Stream Loop (The "Broadcaster")
	ctx := r.Context()
	for {
		select {
		case msg := <-outputChan:
			safeMsg := strings.ReplaceAll(msg, "\n", " ")
			fmt.Fprintf(w, "data: %s\n\n", safeMsg)
			w.(http.Flusher).Flush()
		case <-done:
			return
		case <-ctx.Done():
			// Client disconnected (stop button, refresh, closed tab)
			close(done)
			return
		}
	}
}
