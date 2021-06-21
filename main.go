// A simple client / server to demonstrate using a CA to use mutual certificate
// authentication using HAProxy.
package main

import (
	"flag"
	"fmt"
	"io"
	"log"
	"net/http"
	"time"
)

var mode = flag.String("mode", "server", "Act as `client` or `server`")

// These are set in client-proxy.cfg and server-proxy.cfg respectively
var clientPort = 8080
var serverPort = 3000

func main() {
	flag.Parse()

	if *mode == "server" {
		runServer()
	} else if *mode == "client" {
		runClient()
	} else {
		fmt.Printf("Invalid mode %q; use client or server.", *mode)
	}

}

func runClient() {
	url := fmt.Sprintf("http://localhost:%d", clientPort)
	log.Printf("Starting client: GET %s.", url)

	for {
		time.Sleep(1 * time.Second)

		resp, err := http.Get(url)
		if err != nil {
			log.Printf("Error accessing %s: %v", url, err)
			continue
		}

		data, err := io.ReadAll(resp.Body)
		resp.Body.Close()
		if err != nil {
			log.Printf("Error reading data from %s: %v", url, err)
		}
		log.Printf("Response: %s", string(data))
	}
}

func runServer() {
	log.Printf("Starting server on port %d.", serverPort)

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		log.Print("Received request.")
		fmt.Fprintln(w, "Hello")
	})

	log.Fatal(http.ListenAndServe(
		fmt.Sprintf("localhost:%d", serverPort),
		nil))
}
