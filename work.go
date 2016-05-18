package main

import (
	"bytes"
	"fmt"
	"net/http"
)

type WorkRequests struct {
	Data []WorkRequest `json:"data"`
}

type WorkRequest struct {
	URL     string `json:"url"`
	Cookies string `json:"cookies"`
	Headers string `json:"headers"`
}

func MakeRequest(w WorkRequest) (string, error) {
	var query = []byte("")
	req, err := http.NewRequest("GET", w.URL, bytes.NewBuffer(query))

	req.Header.Set("Content-Type", "text/plain")
	//req.Header.Set()
	//req.AddCookie()

	client := &http.Client{}
	resp, err := client.Do(req)

	if err != nil {
		fmt.Printf("%s %s", req.URL, err)
		return "", err
	}

	defer resp.Body.Close()

	return resp.Status, nil
}
