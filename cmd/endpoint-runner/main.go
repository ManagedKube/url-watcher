package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"
)

func main() {

	log.SetFlags(log.LstdFlags | log.Lshortfile)

	log.Println("[INFO] Starting endpoint-runner.")

	/////////////////////////////
	// Init env params
	/////////////////////////////
	endpointTestJson := os.Getenv("ENDPOINT_TEST_JSON")
	if endpointTestJson == "" {
		fmt.Fprintf(os.Stderr, "ENDPOINT_TEST_JSON environment variable must be set.\n")
		os.Exit(1)
	}


	log.Println("[INFO] ENDPOINT_TEST_JSON: ", endpointTestJson)


	enndpointSpec := []urlWatchEndpointSpec{
		{},
	}
	endpoints := urlWatchEndpoints{
		Endpoints: enndpointSpec,
	}
	urlWatchSpecParsed := urlWatchSpec{
		Watch: endpoints,
	}

	err := json.Unmarshal([]byte(endpointTestJson), &urlWatchSpecParsed)
	if err != nil {
		log.Println(err, "[ERROR] Failed to unmarchall json: ENDPOINT_TEST_JSON")
	}

	log.Println("[INFO] Endpoints to watch:", len(urlWatchSpecParsed.Watch.Endpoints))

	for _, endpoint := range urlWatchSpecParsed.Watch.Endpoints{
		log.Println("[INFO] Endpoint:", endpoint.Host, "| Path:", endpoint.Path, "| Protocol:", endpoint.Protocol)

		// Start a test runner and start testing this endpoint
		go testRunner(endpoint)
	}

	/////////////////////////////
	// http server listen
	/////////////////////////////
	log.Println("[INFO] Server listening")
	http.ListenAndServe(":3000", nil)
}

type urlWatchSpec struct{
	Watch urlWatchEndpoints `json:"watch"`
}

type urlWatchEndpoints struct{
	Endpoints []urlWatchEndpointSpec `json:"endpoints"`
}

type urlWatchEndpointSpec struct{
	Interval int64 `json:"interval,omitempty"`
	Protocol string `json:"protocol,omitempty"`
	Host string `json:"host,omitempty"`
	Method string `json:"method,omitempty"`
	Path string `json:"path,omitempty"`
	Payload string `json:"payload,omitempty"`
	ScrapeTimeout int64 `json:"scrapeTimeout,omitempty"`
}

func testRunner(endpoint urlWatchEndpointSpec){
	log.Println("[INFO] Starting test runner for:", endpoint.Host)

	for true{
		// Sleep
		time.Sleep(time.Duration(endpoint.Interval) * time.Second)

		log.Println("[INFO] running", endpoint.Host)

		testEndpoint(endpoint)
	}

}

func testEndpoint(endpoint urlWatchEndpointSpec){

	switch(endpoint.Method) {
	case "GET":
		log.Println("[INFO] GET")
		runnerGet(endpoint)
	case "POST":
		log.Println("[INFO] POST")
	default:
		log.Println("[ERROR] Didn't find test Method")
	}
}

func runnerGet(endpoint urlWatchEndpointSpec){

	// Doc: https://golang.org/pkg/net/http/
	client := &http.Client{}

	req, err := http.NewRequest("GET", endpoint.Protocol+"://"+endpoint.Host+endpoint.Path, nil)
	if err != nil {
		log.Println("[INFO] Test failed: ", err)
	}

	// Add Optional headers
	// TODO: make this dynamic and take this config in from the ENDPOINT_TEST_JSON input
	req.Header.Add("If-None-Match", `W/"wyzzy"`)

	// Run request
	resp, err := client.Do(req)
	if err != nil {
		log.Println("[INFO] Test failed: ", err)
	}else {

		//Print the HTTP Status Code and Status Name
		log.Println("HTTP Response Status:", endpoint.Host, resp.StatusCode, http.StatusText(resp.StatusCode))

		if resp.StatusCode >= 200 && resp.StatusCode <= 299 {
			log.Println("HTTP Status is in the 2xx range", endpoint.Host)
		} else {
			log.Println("Argh! Broken", endpoint.Host)
		}
	}
}