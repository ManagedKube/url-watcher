package main

import (
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	httpstat "github.com/tcnksm/go-httpstat"
)

var (
	cpuTemp = prometheus.NewGauge(prometheus.GaugeOpts{
		Name: "cpu_temperature_celsius",
		Help: "Current temperature of the CPU.",
	})
	hdFailures = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "hd_errors_total",
			Help: "Number of hard-disk errors.",
		},
		[]string{"device"},
	)

	// Doc: https://godoc.org/github.com/prometheus/client_golang/prometheus
	// Doc: https://godoc.org/github.com/prometheus/client_golang/prometheus#GaugeVec
	promEndpointStatus = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name:      "managedkube_url_watcher_endpoint_status",
			Help:      "The status of an endpoint. alive=1, not reachable=2",
		},
		[]string{
			"endpoint",
			"ingress_name",
			"namespace",
		},
	)


)

func init() {
	// Metrics have to be registered to be exposed:
	prometheus.MustRegister(promEndpointStatus)
}


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
	//log.Println("[INFO] Server listening")
	//http.ListenAndServe(":3000", nil)

	// The Handler function provides a default handler to expose metrics
	// via an HTTP server. "/metrics" is the usual endpoint for that.
	http.Handle("/metrics", promhttp.Handler())
	log.Fatal(http.ListenAndServe(":9093", nil))
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

	// Create a httpstat powered context
	// Doc: https://medium.com/@deeeet/trancing-http-request-latency-in-golang-65b2463f548c
	var result httpstat.Result
	ctx := httpstat.WithHTTPStat(req.Context(), &result)
	req = req.WithContext(ctx)

	// Run request
	resp, err := client.Do(req)
	if err != nil {
		log.Println("[INFO] Test failed: ", err)
	}else {

		if _, err := io.Copy(ioutil.Discard, resp.Body); err != nil {
			log.Fatal(err)
		}
		resp.Body.Close()
		//end := time.Now()

		// Show the httpstat results
		log.Printf("DNS lookup: %d ms", int(result.DNSLookup/time.Millisecond))
		log.Printf("TCP connection: %d ms", int(result.TCPConnection/time.Millisecond))
		log.Printf("TLS handshake: %d ms", int(result.TLSHandshake/time.Millisecond))
		log.Printf("Server processing: %d ms", int(result.ServerProcessing/time.Millisecond))
		log.Printf("Content transfer: %d ms", int(result.ContentTransfer(time.Now())/time.Millisecond))

		//Print the HTTP Status Code and Status Name
		log.Println("HTTP Response Status:", endpoint.Host, resp.StatusCode, http.StatusText(resp.StatusCode))

		if resp.StatusCode >= 200 && resp.StatusCode <= 299 {
			log.Println("HTTP Status is in the 2xx range", endpoint.Host)
		} else {
			log.Println("Argh! Broken", endpoint.Host)
		}
	}

	promEndpointStatus.With(prometheus.Labels{"endpoint": endpoint.Host, "ingress_name": "bar", "namespace": "default"}).Set(float64(result.ServerProcessing/time.Millisecond))
}