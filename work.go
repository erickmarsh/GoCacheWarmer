package main

type WorkRequests struct {
	Data []WorkRequest `json:"data"`
}

type WorkRequest struct {
	URL     string `json:"url"`
	Cookies string `json:"cookies"`
	Headers string `json:"headers"`
}
