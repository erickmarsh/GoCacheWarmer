package main

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"net/http"
)

var url = fmt.Sprintf("http://%s/urls/eric", *HTTPQueue)

func FetchBatchFromRemote(oldBatchID string) (string, []byte, error) {

	return doRequest("GET", url, "")

}

func DeleteBatchFromRemote(oldBatchID string) (string, []byte, error) {
	if oldBatchID == "" {
		return "", []byte(""), nil
	}

	b := "{\"_id\": \"" + oldBatchID + "\"}"

	return doRequest("DELETE", url, b)

}

func doRequest(method string, URL string, reqBody string) (string, []byte, error) {
	client := &http.Client{}

	b := bytes.NewBufferString(reqBody)
	r, err := http.NewRequest(method, URL, b)
	r.Header.Add("Content-Type", "application/json")

	if err != nil {
		fmt.Printf("Error creating request: %s\n", err)
		return "", []byte(""), err
	}

	resp, err := client.Do(r)

	if err != nil {
		fmt.Printf("Error making request: %s\n", err)
		return "", []byte(""), err
	}

	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)

	if err != nil {
		fmt.Printf("Could not parse body: %s\n ", err)
		return "", []byte(""), err
	}

	return resp.Status, body, nil
}
