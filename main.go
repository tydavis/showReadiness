package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"sync"
	"time"
)

// Globals are generally discouraged, but this is a very simple program, and we
// are designing with concurrent access in mind.

// Creating a custom struct with a lock, so we can lock access to the object
type status struct {
	value int
	sync.Mutex
}

// ReadyValue indiciates the program is ready to receive traffic
var ReadyValue = status{value: http.StatusOK}

// LiveValue indiciates the program is alive and should not be terminated
var LiveValue = status{value: http.StatusOK}
var hostname string

func makeNotReady(w http.ResponseWriter, r *http.Request) {
	// Lock access to the variable, then set our global ReadyValue to a failing
	// value before sending the response
	ReadyValue.Lock()
	ReadyValue.value = http.StatusBadRequest
	ReadyValue.Unlock()
	w.Header().Set("responding-pod", hostname)
	fmt.Fprintf(w, "%s", "Set Readiness Value to a failure state")

}

// Hitting this endpoint (by any means) will
func makePodReady(w http.ResponseWriter, r *http.Request) {
	ReadyValue.Lock()
	ReadyValue.value = http.StatusOK
	ReadyValue.Unlock()
	w.Header().Set("responding-pod", hostname)
	fmt.Fprintf(w, "%s", "Set Readiness Value to successful (OK) state")

}

// This function guarantees Kubernetes will kill the pod
func killMe(w http.ResponseWriter, r *http.Request) {
	LiveValue.Lock()
	LiveValue.value = http.StatusBadRequest
	LiveValue.Unlock()
	w.Header().Set("responding-pod", hostname)
	fmt.Fprintf(w, "%s", "Set Liveness Value to a failure state")
}

// Provides Ready state to Kubernetes for receiving/stopping traffic
func readinessCheck(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("responding-pod", hostname)
	ReadyValue.Lock()
	http.Error(w, "Responding with ReadyValue", ReadyValue.value)
	ReadyValue.Unlock()
}

// Endpoint checked internally by Kubernetes
func livenessCheck(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("responding-pod", hostname)
	LiveValue.Lock()
	http.Error(w, "Responding with LiveValue", LiveValue.value)
	LiveValue.Unlock()
}

// The default endpoint
func rootHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("responding-pod", hostname)
	fmt.Fprintf(w, "%s", "I'm serving traffic!")
}

func main() {
	// Force log output to stdout for Docker
	log.SetOutput(os.Stdout)

	// Finish startup
	hostname, _ = os.Hostname()
	log.Println("Service started on port 80")

	// Set up a custom mutex to add custom handlers. [Note] this level of breakout
	// is not required, as we could reasonably compact this into a single call,
	// but would lose some readability.  This also makes it easy for us to use a
	// custom muxer (e.g. gorillamux), custom server, etc
	mux := http.NewServeMux()
	mux.HandleFunc("/", rootHandler)
	mux.HandleFunc("/ping", livenessCheck)
	mux.HandleFunc("/ready", readinessCheck)
	mux.HandleFunc("/makeNotReady", makeNotReady)
	mux.HandleFunc("/makePodReady", makePodReady)
	mux.HandleFunc("/killMe", killMe)

	server := &http.Server{
		Addr:         ":80",
		Handler:      mux,
		ReadTimeout:  3 * time.Second,
		WriteTimeout: 3 * time.Second,
	}

	// Actually run our server
	if err := server.ListenAndServe(); err != nil {
		log.Fatalln(err)
	}
}
