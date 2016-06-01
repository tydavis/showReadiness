package main

import (
	"log"
	"net/http"
	"os"
	"time"

	"github.com/facebookgo/httpdown"
)

var ReadyValue = http.StatusOK
var LiveValue = http.StatusOK
var hostname string

func readinessCheck(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("responding-pod", hostname)
	http.Error(w, "", ReadyValue)
}

func livenessCheck(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("responding-pod", hostname)
	http.Error(w, "", LiveValue)
}

func rootHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("responding-pod", hostname)
	w.Fprintf(w, "%s", "I'm serving traffic!")
}

func main() {
	// Force log output to stdout for Docker
	log.SetOutput(os.Stdout)
	// Delay for startup
	time.Sleep(1 * time.Second)
	hostname, _ = os.Hostname()

	log.Println("Service started on port 80")

	mux := http.NewServeMux()
	mux.HandleFunc("/", rootHandler)
	mux.HandleFunc("/ping", livenessCheck)
	mux.HandleFunc("/ready", readinessCheck)
	mux.HandleFunc("/makeNotReady", makeNotReady)
	mux.HandleFunc("/makePodReady", makePodReady)

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
