package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"encoding/json"
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