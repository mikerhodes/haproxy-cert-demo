// A simple client / server to demonstrate using a CA to use mutual certificate
// authentication using HAProxy. The client and server include reasonable
// default configurations for production use, because several Go defaults
// don't suit production servers.
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

	// This client (somewhat needlessly) uses a reasonable set of defaults for
	// API production use in terms of the c's Timeout and the Transport's
	// connection pool.
	t := http.DefaultTransport.(*http.Transport).Clone()
	t.MaxIdleConns = 100
	t.MaxConnsPerHost = 10
	t.MaxIdleConnsPerHost = 10

	c := &http.Client{
		Timeout:   1 * time.Second,
		Transport: t,
	}

	for {
		time.Sleep(1 * time.Second)

		resp, err := c.Get(url)
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

	// This shows a reasonable default Go server configuration to expose
	// directly to the internet. Primarily it includes timeouts to prevent
	// evil clients hogging all the connections. Likely TLS should be included
	// but isn't here because the point of the example is to show the HAProxy
	// secure connection.
	mux := http.NewServeMux()
	mux.Handle("/", apiHandler{})

	srv := &http.Server{
		Addr:         fmt.Sprintf("localhost:%d", serverPort),
		Handler:      mux,
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 10 * time.Second,
		IdleTimeout:  120 * time.Second,
	}
	log.Println(srv.ListenAndServe())
}

type apiHandler struct{}

func (apiHandler) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	log.Print("Received request.")
	fmt.Fprintln(w, "Hello")
}
