package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
)

var (
	requestNr int64  = 0
	host      string = "unknown"
	isAlive   bool   = true
	isReady   bool   = true
)

func main() {

	finish := make(chan bool)

	currentHost, err := os.Hostname()

	if err != nil {
		log.Println("Could not determine host name:", err)
	} else {
		host = currentHost
	}

	log.Println("Starting Hello Server on " + host)

	server80 := http.NewServeMux()
	server80.HandleFunc("/", Hello)

	go func() {
		err := http.ListenAndServe(":80", server80)
		if err != nil {
			log.Fatal(err)
		}
	}()

	server8080 := http.NewServeMux()
	server8080.HandleFunc("/alive", aliveCheck)
	server8080.HandleFunc("/ready", readyCheck)
	server8080.HandleFunc("/toggleAlive", toggleAlive)
	server8080.HandleFunc("/toggleReady", toggleReady)

	go func() {
		err := http.ListenAndServe(":8080", server8080)
		if err != nil {
			log.Fatal(err)
		}
	}()

	<-finish
}

// Hello gives out a simple hello message
func Hello(w http.ResponseWriter, r *http.Request) {

	requestNr++
	message := fmt.Sprintf("Go Aloha %d from %s on %s ./%s\n", requestNr, host, r.Method, r.URL.Path[1:])
	log.Print(message)
	fmt.Fprint(w, message)
}

func aliveCheck(w http.ResponseWriter, r *http.Request) {

	if isAlive {
		log.Print("Liveness Check: Alive")
		w.WriteHeader(http.StatusOK)
	} else {
		log.Print("Liveness Check: Not Alive")
		w.WriteHeader(http.StatusConflict)
	}
}

func readyCheck(w http.ResponseWriter, r *http.Request) {

	if isReady {
		log.Print("Readiness Check: Ready")
		w.WriteHeader(http.StatusOK)
	} else {
		log.Print("Readiness Check: Not Ready")
		w.WriteHeader(http.StatusConflict)
	}
}

func toggleAlive(w http.ResponseWriter, r *http.Request) {

	isAlive = !isAlive

	w.WriteHeader(http.StatusOK)
}

func toggleReady(w http.ResponseWriter, r *http.Request) {

	isReady = !isReady

	w.WriteHeader(http.StatusOK)
}
