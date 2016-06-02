package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/facebookgo/httpdown"
)

var ReadyValue = http.StatusOK
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

	// Delay for startup
	var delay time.Duration = (1 * time.Second)
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
		Addr:    ":80",
		Handler: mux,
	}

	hd := &httpdown.HTTP{
		StopTimeout: 10 * time.Second,
		KillTimeout: 1 * time.Second,
	}

	if err := httpdown.ListenAndServe(server, hd); err != nil {
		log.Fatalln(err)
	}

}
