package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"strconv"
)

// A buffered channel that we can send work requests on.
var WorkQueue = make(chan WorkRequest, 100)

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

	workItemsAddedCount, err := addWorkToQueue(body, w)

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

func addWorkToQueue(body []byte, w http.ResponseWriter) (int, error) {
	workReqs := WorkRequests{}
	err := json.Unmarshal(body, &workReqs)

	if err != nil {
		return 0, errors.New("Error parsing json")
	}

	for _, work := range workReqs.Data {
		WorkQueue <- work
	}

	return len(workReqs.Data), nil
}
