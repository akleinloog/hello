/*
Copyright © 2020 Arnoud Kleinloog

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/
package cmd

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var (
	requestNr  int64  = 0
	host       string = "unknown"
	serverPort int
	greeting   string
	isAlive    bool = true
	isReady    bool = true
)

// serveCmd represents the serve command
var serveCmd = &cobra.Command{
	Use:   "serve",
	Short: "Starts the HTTP Server.",
	Long: `Starts the HTTP Server that will return a simple hello on any request.
The default port is 80, a different port can be specified if so desired.`,
	Run: func(cmd *cobra.Command, args []string) {
		listen()
	},
}

func init() {
	rootCmd.AddCommand(serveCmd)

	serveCmd.Flags().IntVarP(&serverPort, "port", "p", 80, "port number")
	serveCmd.Flags().StringVarP(&greeting, "greeting", "g", "Hello", "greeting text")

	viper.BindPFlag("port", serveCmd.Flags().Lookup("port"))
	viper.BindPFlag("greeting", serveCmd.Flags().Lookup("greeting"))
}

func Port() int {
	return viper.GetInt("port")
}

func Greeting() string {
	return viper.GetString("greeting")
}

func listen() {

	finish := make(chan bool)

	currentHost, err := os.Hostname()

	if err != nil {
		log.Println("Could not determine host name:", err)
	} else {
		host = currentHost
	}

	log.Printf("Starting %s Server on %s port %d", Greeting(), host, Port())

	server80 := http.NewServeMux()
	server80.HandleFunc("/", hello)

	go func() {
		err := http.ListenAndServe(fmt.Sprintf(":%d", Port()), server80)
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

func hello(w http.ResponseWriter, r *http.Request) {

	requestNr++
	message := fmt.Sprintf("Go %s %d from %s on %s ./%s\n", Greeting(), requestNr, host, r.Method, r.URL.Path[1:])
	log.Print(message)
	fmt.Fprint(w, message)
}

func aliveCheck(w http.ResponseWriter, r *http.Request) {

	if isAlive {
		log.Printf("Liveness Check: Alive on %s\n", host)
		w.WriteHeader(http.StatusOK)
	} else {
		log.Printf("Liveness Check: Not Alive on %s\n", host)
		w.WriteHeader(http.StatusConflict)
	}
}

func readyCheck(w http.ResponseWriter, r *http.Request) {

	if isReady {
		log.Printf("Readiness Check: Ready on %s\n", host)
		w.WriteHeader(http.StatusOK)
	} else {
		log.Printf("Readiness Check: Not Ready on %s\n", host)
		w.WriteHeader(http.StatusConflict)
	}
}

func toggleAlive(w http.ResponseWriter, r *http.Request) {

	log.Printf("Toggle Liveness Check on %s\n", host)

	isAlive = !isAlive

	w.WriteHeader(http.StatusOK)
}

func toggleReady(w http.ResponseWriter, r *http.Request) {

	log.Printf("Toggle Readiness Check on %s\n", host)

	isReady = !isReady

	w.WriteHeader(http.StatusOK)
}
