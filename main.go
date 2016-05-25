package main

import (
	"flag"
	"fmt"
	"net/http"
)

var (
	NWorkers         = flag.Int("n", 4, "The number of workers to start")
	HTTPAddr         = flag.String("http", "127.0.0.1:8002", "Address to listen for HTTP requests on")
	HTTPQueue        = flag.String("q", "127.0.0.1:8001", "Address of URL Queue")
	PostRequestPause = flag.Int64("p", 1, "The number of seconds to pause after a request")
)

func main() {
	// Parse the command-line flags.
	flag.Parse()

	// Start the dispatcher.
	fmt.Println("Starting the dispatcher")
	//StartDispatcher(*NWorkers)

	// Register our collector as an HTTP handler function.
	fmt.Println("Registering the collector")
	http.HandleFunc("/work", Collector)
	http.HandleFunc("/workers/", Workers)

	fmt.Println("Starting batch handler")
	go BatchInitialize()

	// Start the HTTP server!
	fmt.Println("HTTP server listening on", *HTTPAddr)

	if err := http.ListenAndServe(*HTTPAddr, nil); err != nil {
		fmt.Println(err.Error())
	}

}
