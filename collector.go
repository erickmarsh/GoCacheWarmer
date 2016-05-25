package main

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"strconv"
)

// WorkQueue - A buffered channel that we can send work requests on.
var WorkQueue = make(chan WorkRequest, 100)

// ActiveWorkRequests - a map of the active url batches
//var ActiveWorkRequests = make(map[string]WorkRequests)

// Collector is the main handler for accepting new urls via a POST
func Collector(w http.ResponseWriter, r *http.Request) {
	// Make sure we can only be called with an HTTP POST request.
	if r.Method != "POST" {
		w.Header().Set("Allow", "POST")
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	defer r.Body.Close()
	body, err := ioutil.ReadAll(r.Body)

	if err != nil {
		http.Error(w, "No JSON Body Posted", http.StatusBadRequest)
		return
	}

	workItemsAddedCount, err := addWorkToQueue(body)

	if err != nil {
		http.Error(w, "Cannot parse json body", http.StatusBadRequest)
		return
	}

	// Push the work onto the queue.
	// WorkQueue <- work
	fmt.Println(strconv.Itoa(workItemsAddedCount) + " Work request queued")

	// And let the user know their work request was created.
	w.WriteHeader(http.StatusCreated)
	return
}
