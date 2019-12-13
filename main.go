package main

import (
	"encoding/json"
	"fmt"
	"github.com/DataDog/datadog-go/statsd"
	"io/ioutil"
	"log"
	"net/http"
)

type Data struct {
	//TransID int64 `json:"TRANSID"`
	// RESPONSE       string `json:"RESPONSE"`
	// EMAIL string `json:"EMAIL"`
	// TIMESTAMP      int    `json:"TIMESTAMP"`
	// FROMADDRESS    string `json:"FROMADDRESS"`
	Event string `json:"EVENT"`
	// MSIZE          int    `json:"MSIZE"`
	// USERAGENT      string `json:"USERAGENT"`
	// TAGS           string `json:"TAGS"`
	// XAPIHEADER     string `json:"X-APIHEADER"`
	// URL            string `json:"URL"`
	// IPADDRESS      string `json:"IPADDRESS"`
	// BOUNCETYPE     string `json:"BOUNCE_TYPE"`
	// BOUNCEREASON   string `json:"BOUNCE_REASON"`
	// BOUNCEREASONID int    `json:"BOUNCE_REASONID"`
}

var statsdClient *statsd.Client

func handler(w http.ResponseWriter, r *http.Request) {
	var data []Data
	body, err := ioutil.ReadAll(r.Body)
	check(err)
	err = json.Unmarshal(body, &data)
	if err == nil {
		for _, d := range data {
			metric := fmt.Sprintf("pepipost.email.%s", d.Event)
			err = statsdClient.Incr(metric, nil, 1)
			check(err)
		}
	}
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("OK"))
}
func healthHandler(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("OK"))
}

func readinessHandler(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("OK"))
}

func check(err error) {
	if err != nil {
		log.Fatal(err)
	}
}
func main() {

	cl, err := statsd.New("")
	check(err)
	statsdClient = cl

	log.Println("Server Started")

	http.HandleFunc("/", handler)
	http.HandleFunc("/healthz", healthHandler)
	http.HandleFunc("/readiness", readinessHandler)
	log.Fatal(http.ListenAndServe(":8080", nil))
}
