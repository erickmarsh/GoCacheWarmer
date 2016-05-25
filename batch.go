package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"sync"
	"time"
)

var BatchCompleteChannel = make(chan string, 5)
var BatchActiveChannel = make(chan string, 5)

type Batch struct {
	ID      string
	retries int32
}

var BatchWaitGroup sync.WaitGroup

func BatchInitialize() {

	fmt.Println("Starting the Batch Life Cycle")
	//BatchCompleteChannel <- ""

	go startBatchLifecycle()

}

func startBatchLifecycle() {

	go StartDispatcher(*NWorkers)
	// Since this is the first run there is no lastBatchID

	var lastBatchID string

	for {

		fmt.Println("Getting Next Batch")

		go FetchNextBatch(lastBatchID)
		lastBatchID = <-BatchActiveChannel

		fmt.Println(string(lastBatchID))
		BatchWaitGroup.Wait()
		fmt.Println("Batch is done - Deleting " + string(lastBatchID))
		BatchComplete(lastBatchID)
	}

}

// FetchNextBatch is what is done to get a new batch of URLs
func FetchNextBatch(lastBatchID string) {
	fmt.Println("\n\nFetching next batch of URLS")

	hasData := false
	var body []byte

	for hasData == false {
		var status string
		var err error

		status, body, err = FetchBatchFromRemote(lastBatchID)

		if err != nil {
			fmt.Printf("Didn't parse new batch %s\n", err)
			fmt.Printf("Retrying...")
			BatchCompleteChannel <- ""
			return
		}

		if status == "200 OK" { // && (string(body) != "null" && string(body) != "") {
			hasData = true
		} else {
			fmt.Println("Sleeping....")
			time.Sleep(5 * time.Second)
		}

	}

	queueItemsAdded, err := addWorkToQueue(body)

	if err != nil {
		fmt.Printf("Did not add work to queue: %s\n", err)
		fmt.Printf("%s\n", body)
	}

	fmt.Printf("Added %d items to work queue\n", queueItemsAdded)
}

// BatchComplete is the actions taken when all urls in a batch are
// executed
func BatchComplete(lastBatchID string) {
	fmt.Println("Deleting Batch: " + lastBatchID)

	status, body, err := DeleteBatchFromRemote(lastBatchID)

	if err != nil {
		fmt.Println("Failed to delete from remote " + lastBatchID)
		fmt.Println("Not sure what to do now")
		return
	}

	if status != "200 OK" {
		fmt.Printf("HTTP Error on batch delete: %s \n%s", status, body)
		return
	}

	fmt.Println("Batch Deleted")

}

func addWorkToQueue(body []byte) (int, error) {
	workReqs := WorkRequests{}

	err := json.Unmarshal(body, &workReqs)

	if err != nil {
		return 0, err
	}

	BatchActiveChannel <- workReqs.ID

	fmt.Print("\n\n--------------------------------------\n")
	fmt.Printf("Adding %d items to work queue\n", len(workReqs.URLs))

	for _, work := range workReqs.URLs {
		fmt.Println("adding to wait group")
		BatchWaitGroup.Add(1)
		// add the ID to the work item so each knows which batch it is part of
		work.BatchID = workReqs.ID
		WorkQueue <- work
	}

	//ActiveWorkRequests[workReqs.ID] = workReqs

	return len(workReqs.URLs), nil
}

/*


   OLD Code


*/

func deleteOldURLBatch(oldBatchID string) {
	if oldBatchID == "" {
		return
	}

	client := &http.Client{}

	b := bytes.NewBufferString("{\"_id\": \"" + oldBatchID + "\"}")
	r, _ := http.NewRequest("DELETE", "http://"+*HTTPQueue+"/urls/eric", b)
	r.Header.Add("Content-Type", "application/json")

	resp, _ := client.Do(r)
	fmt.Println(resp.Status)

}

func getNewURLBatch(oldBatchID string) (string, []byte, error) {

	client := &http.Client{}
	b := bytes.NewBufferString("")
	r, _ := http.NewRequest("GET", "http://"+*HTTPQueue+"/urls/eric", b)
	r.Header.Add("Content-Type", "application/json")

	resp, _ := client.Do(r)
	fmt.Println(resp.Status)

	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)

	if err != nil {
		//http.Error(w, "No JSON Body Posted", http.StatusBadRequest)
		fmt.Println("No JSON Body Posted")
		return "", nil, err
	}

	fmt.Println(resp.Status)
	fmt.Printf("%+v", string(body))

	return resp.Status, body, nil

}
