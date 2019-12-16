package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"

	"github.com/DataDog/datadog-go/statsd"
	log "github.com/sirupsen/logrus"
)

type Data struct {
	TransID        int64  `json:"TRANSID"`
	Response       string `json:"RESPONSE"`
	Email          string `json:"EMAIL"`
	Timestamp      int    `json:"TIMESTAMP"`
	FromAddress    string `json:"FROMADDRESS"`
	Event          string `json:"EVENT"`
	MSize          int    `json:"MSIZE"`
	UserAgent      string `json:"USERAGENT"`
	Tags           string `json:"TAGS"`
	XAPIHeader     string `json:"X-APIHEADER"`
	URL            string `json:"URL"`
	IPAddress      string `json:"IPADDRESS"`
	BounceType     string `json:"BOUNCE_TYPE"`
	BounceReason   string `json:"BOUNCE_REASON"`
	BounceReasonID int    `json:"BOUNCE_REASONID"`
}

func init() {
	log.SetFormatter(&log.JSONFormatter{})
	log.SetOutput(os.Stdout)
	log.SetLevel(log.WarnLevel)
}

func main() {
	svc, err := statsd.New("")

	if err != nil {
		log.Fatalln(err)
	}

	log.Info("Server Started")

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		var data []Data

		body, err := ioutil.ReadAll(r.Body)

		if err != nil {
			log.Fatalln(err)
		}

		err = json.Unmarshal(body, &data)

		if err == nil {
			for _, d := range data {
				metric := fmt.Sprintf("pepipost.email.%s", d.Event)
				err = svc.Incr(metric, nil, 1)

				if d.Event == "invalid" {
					log.WithField("event", d).Error("pepipost.email.invalid")
				}
			}
		}

		w.WriteHeader(http.StatusOK)
		w.Write([]byte("OK"))
	})

	http.HandleFunc("/healthz", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("OK"))
	})

	http.HandleFunc("/readiness", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("OK"))
	})

	err = http.ListenAndServe(":8080", nil)

	if err != nil {
		log.Fatalf("Could not start server: %s\n", err.Error())
	}
}
