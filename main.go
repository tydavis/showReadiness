package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"time"
)

// ReadyValue indiciates the program is ready to receive traffic
var ReadyValue = http.StatusOK

// LiveValue indiciates the program is alive and should not be terminated
var LiveValue = http.StatusOK
var hostname string

func makeNotReady(w http.ResponseWriter, r *http.Request) {
	ReadyValue = http.StatusBadRequest
	w.Header().Set("responding-pod", hostname)
	fmt.Fprintf(w, "%s", "Set Readiness Value to a failure state")

}

func makePodReady(w http.ResponseWriter, r *http.Request) {
	ReadyValue = http.StatusOK
	w.Header().Set("responding-pod", hostname)
	fmt.Fprintf(w, "%s", "Set Readiness Value to successful (OK) state")

}

func killMe(w http.ResponseWriter, r *http.Request) {
	LiveValue = http.StatusBadRequest
	w.Header().Set("responding-pod", hostname)
	fmt.Fprintf(w, "%s", "Set Liveness Value to a failure state")
}

func readinessCheck(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("responding-pod", hostname)
	http.Error(w, "Responding with ReadyValue", ReadyValue)
}

func livenessCheck(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("responding-pod", hostname)
	http.Error(w, "Responding with LiveValue", LiveValue)
}

func rootHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("responding-pod", hostname)
	fmt.Fprintf(w, "%s", "I'm serving traffic!")
}

func main() {
	// Force log output to stdout for Docker
	log.SetOutput(os.Stdout)

	// Configurable delay for startup
	var delay = (1 * time.Second)
	if os.Getenv("APPDELAY") != "" {
		var err error
		delay, err = time.ParseDuration(os.Getenv("APPDELAY"))
		if err != nil {
			log.Fatalf("Failed to parse time duration: %v", err)
		}
	}
	time.Sleep(delay)

	// Finish startup
	hostname, _ = os.Hostname()
	log.Println("Service started on port 80")

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

	if err := server.ListenAndServe(); err != nil {
		log.Fatalln(err)
	}

}
