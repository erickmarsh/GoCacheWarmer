package main

import (
	"bytes"
	"fmt"
	"net/http"
)

type WorkRequests struct {
	ID          string                 `json:"_id"`
	CallbackURL string                 `json:"callback_url"`
	ClientEnv   string                 `json:"client_env"`
	Priority    int64                  `json:"priority"`
	URLs        map[string]WorkRequest `json:"urls"`
}

type WorkRequest struct {
	BatchID string
	URL     string     `json:"url"`
	IP      string     `json:"ip"`
	Host    string     `json:"host"`
	Cookies []KeyValue `json:"cookies"`
	Headers []KeyValue `json:"headers"`
}

type KeyValue struct {
	Key   string `json:"key"`
	Value string `json:"value"`
}

func MakeRequest(w WorkRequest) (string, error) {
	var query = []byte("")
	req, err := http.NewRequest("GET", w.URL, bytes.NewBuffer(query))

	// add the cookies
	for _, cookie := range w.Cookies {
		c := http.Cookie{Name: cookie.Key, Value: cookie.Value}
		req.AddCookie(&c)
	}

	// add the headers
	for _, header := range w.Headers {
		req.Header.Add(header.Key, header.Value)
	}

	client := &http.Client{}
	resp, err := client.Do(req)

	if err != nil {
		fmt.Printf("%s %s", req.URL, err)
		return "", err
	}

	defer resp.Body.Close()

	return resp.Status, nil
}
