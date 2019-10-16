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
	promEndpointStatusCode = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name:      "managedkube_url_watcher_endpoint_status_code",
			Help:      "The status code of the returned HTTP call",
		},
		[]string{
			"endpoint",
			"ingress_name",
			"namespace",
			"path",
		},
	)
	promEndpointStatus = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name:      "managedkube_url_watcher_endpoint_status",
			Help:      "The status code string of the returned HTTP call",
		},
		[]string{
			"endpoint",
			"ingress_name",
			"namespace",
			"path",
			"status",
		},
	)
	promEndpointProto = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name:      "managedkube_url_watcher_endpoint_proto",
			Help:      "The protocol used",
		},
		[]string{
			"endpoint",
			"ingress_name",
			"namespace",
			"path",
			"proto",
		},
	)
	promEndpointHttpStatsDnsLookupTime = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name:      "managedkube_url_watcher_endpoint_http_stats_dns_lookup_time_ms",
			Help:      "The time taken to look up the DNS",
		},
		[]string{
			"endpoint",
			"ingress_name",
			"namespace",
			"path",
		},
	)
	promEndpointHttpStatsTcpConnectionTime = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name:      "managedkube_url_watcher_endpoint_http_stats_tcp_connection_time_ms",
			Help:      "The time taken for the TCP connection",
		},
		[]string{
			"endpoint",
			"ingress_name",
			"namespace",
			"path",
		},
	)
	promEndpointHttpStatsTlsHandshakeTime = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name:      "managedkube_url_watcher_endpoint_http_stats_tls_handshake_time_ms",
			Help:      "The time taken for the TLS handshake",
		},
		[]string{
			"endpoint",
			"ingress_name",
			"namespace",
			"path",
		},
	)
	promEndpointHttpStatsServerProcessingTime = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name:      "managedkube_url_watcher_endpoint_http_stats_server_processing_time_ms",
			Help:      "The time taken for the server to process the request",
		},
		[]string{
			"endpoint",
			"ingress_name",
			"namespace",
			"path",
		},
	)
	promEndpointHttpStatsContentTransferTime = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name:      "managedkube_url_watcher_endpoint_http_stats_content_transfer_time_ms",
			Help:      "The time taken for to transfer the content",
		},
		[]string{
			"endpoint",
			"ingress_name",
			"namespace",
			"path",
		},
	)

)

func init() {
	// Metrics have to be registered to be exposed:
	prometheus.MustRegister(promEndpointStatusCode)
	prometheus.MustRegister(promEndpointStatus)
	prometheus.MustRegister(promEndpointProto)
	prometheus.MustRegister(promEndpointHttpStatsDnsLookupTime)
	prometheus.MustRegister(promEndpointHttpStatsTcpConnectionTime)
	prometheus.MustRegister(promEndpointHttpStatsTlsHandshakeTime)
	prometheus.MustRegister(promEndpointHttpStatsServerProcessingTime)
	prometheus.MustRegister(promEndpointHttpStatsContentTransferTime)
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


	endpoints := []urlWatchEndpointData{}

	endpointList := urlWatchEndpoints{
		Endpoints: endpoints,
	}

	urlWatchSpecParsed := urlWatchSpec{
		Watch: endpointList,
	}

	err := json.Unmarshal([]byte(endpointTestJson), &urlWatchSpecParsed)
	if err != nil {
		log.Println(err, "[ERROR] Failed to unmarchall json: ENDPOINT_TEST_JSON")
	}

	log.Println("[INFO] Endpoints to watch:", len(urlWatchSpecParsed.Watch.Endpoints))

	for _, endpointData := range urlWatchSpecParsed.Watch.Endpoints{
		if(areParametersOk(endpointData)){
			log.Println("[INFO] Endpoint:", endpointData.Endpoint.Host, "| Path:", endpointData.Endpoint.Path, "| Protocol:", endpointData.Endpoint.Protocol)

			// Start a test runner and start testing this endpoint
			go runner(endpointData)
		}else{
			log.Println("[WARNING] areParametersOk: false, Endpoint:", endpointData.Endpoint.Host, "| Path:", endpointData.Endpoint.Path, "| Protocol:", endpointData.Endpoint.Protocol)

			// Output this endpoint info to prometheus metrics?
			// Probably a good idea so people would know what this was not able to test
		}

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
	Endpoints []urlWatchEndpointData `json:"endpoints"`
}

type urlWatchEndpointData struct{
	MetaData  urlWatchEndpointMetaData `json:"metadata"`
	Endpoint urlWatchEndpointSpec `json:"endpoint"`
}

type urlWatchEndpointMetaData struct{
	Name string `json:"name"`
	Namespace string `json:"namespace"`
	IngressName string `json:"ingressName"`
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

func runner(endpoint urlWatchEndpointData){
	log.Println("[INFO] Starting test runner for:", endpoint.Endpoint.Host)

	for true{
		// Sleep
		time.Sleep(time.Duration(endpoint.Endpoint.Interval) * time.Second)

		log.Println("[INFO] running", endpoint.Endpoint.Host)

		// Run the endpoint watch action
		results := runEndpoint(endpoint)

		// Post prometheus with the results
		updatePrometheusMetrics(results)
	}

}

func runEndpoint(endpoint urlWatchEndpointData) endpointResults{

	var results endpointResults

	switch(endpoint.Endpoint.Method) {
	case "GET":
		log.Println("[INFO] GET")
		results = actionGet(endpoint)
	case "POST":
		log.Println("[INFO] POST")
	default:
		log.Println("[ERROR] Didn't find test Method")
	}

	return results
}

func actionGet(endpoint urlWatchEndpointData) endpointResults{

	var results endpointResults

	// Doc: https://golang.org/pkg/net/http/
	client := &http.Client{}

	req, err := http.NewRequest("GET", endpoint.Endpoint.Protocol+"://"+endpoint.Endpoint.Host+endpoint.Endpoint.Path, nil)
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

		results = endpointResults{
			httpResults{
				StatusCode: -1,
				Status: "Failed",
				Proto: "",
			},
			httpStats{
				DnsLookupTime: 0,
				TcpConnectionTime: 0,
				TlsHandshakeTime: 0,
				ServerProcessingTime: 0,
				ContentTransferTime: 0,
			},
			endpoint,
		}
	}else {

		if _, err := io.Copy(ioutil.Discard, resp.Body); err != nil {
			//log.Fatal(err)
			log.Println("[INFO] Test failed: ", err)

			results = endpointResults{
				httpResults{
					StatusCode: -1,
					Status: "Failed",
					Proto: "",
				},
				httpStats{
					DnsLookupTime: 0,
					TcpConnectionTime: 0,
					TlsHandshakeTime: 0,
					ServerProcessingTime: 0,
					ContentTransferTime: 0,
				},
				endpoint,
			}

		}

		//end := time.Now()

		// Show the httpstat results
		log.Printf("DNS lookup: %d ms", int(result.DNSLookup/time.Millisecond))
		log.Printf("TCP connection: %d ms", int(result.TCPConnection/time.Millisecond))
		log.Printf("TLS handshake: %d ms", int(result.TLSHandshake/time.Millisecond))
		log.Printf("Server processing: %d ms", int(result.ServerProcessing/time.Millisecond))
		log.Printf("Content transfer: %d ms", int(result.ContentTransfer(time.Now())/time.Millisecond))

		//Print the HTTP Status Code and Status Name
		log.Println("HTTP Response Status:", endpoint.Endpoint.Host, resp.StatusCode, http.StatusText(resp.StatusCode))

		if resp.StatusCode >= 200 && resp.StatusCode <= 299 {
			log.Println("HTTP Status is in the 2xx range", endpoint.Endpoint.Host)
		} else {
			log.Println("Argh! Broken", endpoint.Endpoint.Host)
		}

		//promEndpointStatusCode.With(prometheus.Labels{"endpoint": endpoint.Endpoint.Host, "ingress_name": endpoint.MetaData.IngressName, "namespace": endpoint.MetaData.Namespace, "path": endpoint.Endpoint.Path}).Set(float64(resp.StatusCode))
		//promEndpointStatusCode.With(prometheus.Labels{"endpoint": endpoint.Endpoint.Host, "ingress_name": endpoint.MetaData.IngressName, "namespace": endpoint.MetaData.Namespace, "path": endpoint.Endpoint.Path}).Set(float64(result.ServerProcessing/time.Millisecond))

		results = endpointResults{
			httpResults{
				StatusCode: resp.StatusCode,
				Status: resp.Status,
				Proto: resp.Proto,
			},
			httpStats{
				DnsLookupTime: int(result.DNSLookup/time.Millisecond),
				TcpConnectionTime: int(result.TCPConnection/time.Millisecond),
				TlsHandshakeTime: int(result.TLSHandshake/time.Millisecond),
				ServerProcessingTime: int(result.ServerProcessing/time.Millisecond),
				ContentTransferTime: int(result.ContentTransfer(time.Now())/time.Millisecond),
			},
			endpoint,
		}

		resp.Body.Close()
	}

	return results
}

type endpointResults struct{
	Http httpResults `json:"http"`
	HttpStats httpStats `json:"httpStats"`
	EndpointData urlWatchEndpointData `json:"endpointData"`
}

type httpResults struct {
	StatusCode int    `json:"statusCode"`
	Status     string `json:"status"`
	Proto      string `json:"proto"`
}

type httpStats struct {
	DnsLookupTime int `json:"dnsLookupTime"`
	TcpConnectionTime int `json:"tcpConnectionTime"`
	TlsHandshakeTime int `json:"tlsHandshakeTime"`
	ServerProcessingTime int `json:"serverProcessingTime"`
	ContentTransferTime int `json:"contentTransferTime"`
}

func areParametersOk(endpoint urlWatchEndpointData) bool{
	isOk := true

	if(!isValidHostname(endpoint.Endpoint.Host)){
		isOk = false
	}

	return isOk
}

func isValidHostname(hostname string) bool{
	isValid := true

	if (hostname == "") {
		isValid = false
	}

	return isValid
}

func updatePrometheusMetrics(results endpointResults){
	log.Println("[INFO] Updating Prometheus metrics")

	promEndpointStatusCode.With(prometheus.Labels{"endpoint": results.EndpointData.Endpoint.Host, "ingress_name": results.EndpointData.MetaData.IngressName, "namespace": results.EndpointData.MetaData.Namespace, "path": results.EndpointData.Endpoint.Path}).Set(float64(results.Http.StatusCode))
	promEndpointStatus.With(prometheus.Labels{"endpoint": results.EndpointData.Endpoint.Host, "ingress_name": results.EndpointData.MetaData.IngressName, "namespace": results.EndpointData.MetaData.Namespace, "path": results.EndpointData.Endpoint.Path, "status": results.Http.Status}).Set(float64(1))
	promEndpointProto.With(prometheus.Labels{"endpoint": results.EndpointData.Endpoint.Host, "ingress_name": results.EndpointData.MetaData.IngressName, "namespace": results.EndpointData.MetaData.Namespace, "path": results.EndpointData.Endpoint.Path, "proto": results.Http.Proto}).Set(float64(1))

	promEndpointHttpStatsDnsLookupTime.With(prometheus.Labels{"endpoint": results.EndpointData.Endpoint.Host, "ingress_name": results.EndpointData.MetaData.IngressName, "namespace": results.EndpointData.MetaData.Namespace, "path": results.EndpointData.Endpoint.Path}).Set(float64(results.HttpStats.DnsLookupTime))
	promEndpointHttpStatsTcpConnectionTime.With(prometheus.Labels{"endpoint": results.EndpointData.Endpoint.Host, "ingress_name": results.EndpointData.MetaData.IngressName, "namespace": results.EndpointData.MetaData.Namespace, "path": results.EndpointData.Endpoint.Path}).Set(float64(results.HttpStats.TcpConnectionTime))
	promEndpointHttpStatsTlsHandshakeTime.With(prometheus.Labels{"endpoint": results.EndpointData.Endpoint.Host, "ingress_name": results.EndpointData.MetaData.IngressName, "namespace": results.EndpointData.MetaData.Namespace, "path": results.EndpointData.Endpoint.Path}).Set(float64(results.HttpStats.TlsHandshakeTime))
	promEndpointHttpStatsServerProcessingTime.With(prometheus.Labels{"endpoint": results.EndpointData.Endpoint.Host, "ingress_name": results.EndpointData.MetaData.IngressName, "namespace": results.EndpointData.MetaData.Namespace, "path": results.EndpointData.Endpoint.Path}).Set(float64(results.HttpStats.ServerProcessingTime))
	promEndpointHttpStatsContentTransferTime.With(prometheus.Labels{"endpoint": results.EndpointData.Endpoint.Host, "ingress_name": results.EndpointData.MetaData.IngressName, "namespace": results.EndpointData.MetaData.Namespace, "path": results.EndpointData.Endpoint.Path}).Set(float64(results.HttpStats.ContentTransferTime))
}

