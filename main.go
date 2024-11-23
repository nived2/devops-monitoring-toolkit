package main

import (
	"fmt"
	"log"
	"net/http"
	"runtime"
	"time"
)

var memoryLeakSlice []string

func main() {
	// Configure logging
	log.SetFlags(log.LstdFlags | log.Lshortfile)

	// Register handlers
	http.HandleFunc("/", memoryLeakHandler)
	http.HandleFunc("/slow", slowResponseHandler)
	http.HandleFunc("/error", errorHandler)
	http.HandleFunc("/metrics", metricsHandler)
	http.HandleFunc("/health", healthCheckHandler)

	// Start server
	log.Println("Server starting on :8080")
	if err := http.ListenAndServe(":8080", nil); err != nil {
		log.Fatal(err)
	}
}

func memoryLeakHandler(w http.ResponseWriter, r *http.Request) {
	// Simulate memory leak by storing data without cleanup
	memoryLeakSlice = append(memoryLeakSlice, string(make([]byte, 1024*1024))) // Allocate 1MB
	
	var m runtime.MemStats
	runtime.ReadMemStats(&m)
	
	fmt.Fprintf(w, "Memory usage: %v MB\n", m.Alloc/1024/1024)
}

func slowResponseHandler(w http.ResponseWriter, r *http.Request) {
	log.Println("Slow request started")
	time.Sleep(10 * time.Second)
	fmt.Fprintf(w, "Slow response after 10 seconds\n")
	log.Println("Slow request completed")
}

func errorHandler(w http.ResponseWriter, r *http.Request) {
	log.Println("Error endpoint accessed")
	http.Error(w, "Internal Server Error", http.StatusInternalServerError)
}

func metricsHandler(w http.ResponseWriter, r *http.Request) {
	var m runtime.MemStats
	runtime.ReadMemStats(&m)
	
	metrics := fmt.Sprintf(`
# HELP go_memory_alloc_bytes Memory usage in bytes
go_memory_alloc_bytes %v
# HELP go_goroutines Number of goroutines
go_goroutines %v
`, m.Alloc, runtime.NumGoroutine())
	
	w.Header().Set("Content-Type", "text/plain")
	fmt.Fprint(w, metrics)
}

func healthCheckHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	fmt.Fprint(w, `{"status": "healthy"}`)
}
