package main

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strconv"
	"sync"
	"testing"
)

func TestSendAndDrop(t *testing.T) {
	// Create a test server.
	numHits := 0
	numHitsMutex := new(sync.Mutex)
	wg := new(sync.WaitGroup)
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		numHitsMutex.Lock()
		defer wg.Done()
		defer numHitsMutex.Unlock()
		numHits++
	}))
	srvURL, _ := url.Parse(srv.URL)

	// Set the host and port variables.
	*port, _ = strconv.Atoi(srvURL.Port())
	*host = srvURL.Hostname()

	// Set the number of requests to send.
	numRequests := 10
	*num = numRequests
	wg.Add(numRequests)

	// Send the requests.
	sendRequests()
	wg.Wait()

	if numHits != numRequests {
		fmt.Printf("Did not receive number of expected requests. Got %v received %v\n", numHits, numRequests)
		t.Fail()
	}
}
