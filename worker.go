package main

import (
	"fmt"
	"net/http"
	"strconv"
)

// NewWorker creates, and returns a new Worker object. Its only argument
// is a channel that the worker can add itself to whenever it is done its
// work.
func NewWorker(id int, workerQueue chan chan WorkRequest) Worker {
	// Create, and return the worker.
	worker := Worker{
		ID:          id,
		Work:        make(chan WorkRequest),
		WorkerQueue: workerQueue,
		QuitChan:    make(chan bool)}

	return worker
}

type Worker struct {
	ID          int
	Work        chan WorkRequest
	WorkerQueue chan chan WorkRequest
	QuitChan    chan bool
}

// This function "starts" the worker by starting a goroutine, that is
// an infinite "for-select" loop.
func (w *Worker) Start() {
	go func() {
		for {
			// Add ourselves into the worker queue.
			w.WorkerQueue <- w.Work

			select {
			case work := <-w.Work:
				// Receive a work request.
				fmt.Printf("worker%d: Received work request %s\n", w.ID, work.URL)

				status, err := MakeRequest(work)

				if err != nil {
					fmt.Printf("%s: %s\n", work.URL, err.Error())
				}

				fmt.Printf("worker%d: %s: %s\n", w.ID, status, work.URL)

			case <-w.QuitChan:
				// We have been asked to stop.
				fmt.Printf("worker%d stopping\n", w.ID)
				return
			}
		}
	}()
}

// Stop tells the worker to stop listening for work requests.
//
// Note that the worker will only stop *after* it has finished its work.
func (w *Worker) Stop() {
	go func() {
		w.QuitChan <- true
	}()
}

func Workers(w http.ResponseWriter, r *http.Request) {
	// Make sure we can only be called with an HTTP POST request.
	if r.Method != "POST" {
		w.Header().Set("Allow", "POST")
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	//	body, err := ioutil.ReadAll(r.Body)
	//	defer r.Body.Close()
	/*
		if err != nil {
			http.Error(w, "No JSON Body Posted", http.StatusBadRequest)
			return
		}
	*/
	workerCount := r.FormValue("workers")

	i, err := strconv.Atoi(workerCount)

	if err != nil {
		fmt.Println(err)
		http.Error(w, err.Error(), http.StatusBadRequest)
	} else if i > 0 {
		fmt.Println("Adjusting worker count")
		AdjustWorkers(i)
		w.WriteHeader(http.StatusOK)
	} else {
		http.Error(w, "Worker Count must be positive", http.StatusBadRequest)
	}

	return
}
