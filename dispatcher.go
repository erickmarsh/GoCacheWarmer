package main

import "fmt"

var WorkerQueue chan chan WorkRequest

var ActiveWorkers = -1

func StartDispatcher(nworkers int) {

	ActiveWorkers = nworkers

	// First, initialize the channel we are going to but the workers' work channels into.
	WorkerQueue = make(chan chan WorkRequest, nworkers)

	// Now, create all of our workers.
	for i := 0; i < nworkers; i++ {
		fmt.Println("Starting worker", i+1)
		worker := NewWorker(i+1, WorkerQueue)
		worker.Start()
	}

	go func() {
		for {
			select {
			case work := <-WorkQueue:
				fmt.Println("Received work requeust %s", string(work.URL))
				go func() {
					worker := <-WorkerQueue

					//fmt.Println("Dispatching work request")
					worker <- work
				}()
			}
		}
	}()
}

func AdjustWorkers(nworkers int) {

	//workerDiff := nworkers - ActiveWorkers

	fmt.Println("Stopping all workers")

	worker := NewWorker(ActiveWorkers+1, WorkerQueue)

	for i := 0; i < ActiveWorkers; i++ {
		worker.Stop()
	}

	fmt.Printf("Starting %d workers\n", nworkers)

	// First, initialize the channel we are going to but the workers' work channels into.
	WorkerQueue = make(chan chan WorkRequest, nworkers)

	// Now, create all of our workers.
	for i := 0; i < nworkers; i++ {
		fmt.Println("Starting worker", i+1)
		worker := NewWorker(i+1, WorkerQueue)
		worker.Start()
	}

	ActiveWorkers = nworkers

}
