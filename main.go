package main

import (
	"encoding/json"
	"log"
	"net/http"
	"os"

	"github.com/gorilla/mux"
	"github.com/hajhatten/buzzwords"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

var (
	buzzwordCounter = prometheus.NewCounter(prometheus.CounterOpts{
		Name: "buzzword_api",
		Help: "Number of times buzzwords has been returned",
	})
	verbsCounter = prometheus.NewCounter(prometheus.CounterOpts{
		Name: "buzzword_verbs_api",
		Help: "Number of times buzzwords with verbs has been returned",
	})
	suffixCounter = prometheus.NewCounter(prometheus.CounterOpts{
		Name: "buzzword_suffix_api",
		Help: "Number of times buzzwords with suffix has been returned",
	})
	verbsAndSuffixCounter = prometheus.NewCounter(prometheus.CounterOpts{
		Name: "buzzword_verbsuffix_api",
		Help: "Number of times buzzwords with verbs and suffix has been returned",
	})
)

func init() {

	// Initialize prometheus counters
	prometheus.MustRegister(
		buzzwordCounter,
		verbsCounter,
		suffixCounter,
		verbsAndSuffixCounter,
	)
}

func main() {

	r := mux.NewRouter()

	// Api routes
	r.HandleFunc("/buzzword", buzzwordsHandler)
	r.HandleFunc("/suffix", suffixHandler)
	r.HandleFunc("/verb", verbsHandler)
	r.HandleFunc("/verbsuffix", verbsAndSuffixHandler)

	// Prometheus metrics
	r.Handle("/metrics", promhttp.Handler())

	// Static files
	fs := http.FileServer(assetFS())
	r.PathPrefix("/").Handler(fs)

	// Set port for http server
	var port string
	if os.Getenv("PORT") != "" {
		port = ":" + os.Getenv("PORT")
	} else {
		port = ":3001"
	}

	log.Fatal(http.ListenAndServe(port, r))

}
func responseWithJSON(w http.ResponseWriter, body interface{}, code int) {
	result, err := json.Marshal(body)
	if err != nil {
		log.Fatal(err)
	}
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(code)
	w.Write(result)
}

func responseWithText(w http.ResponseWriter, body string, code int) {
	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	w.WriteHeader(code)
	w.Write([]byte(body))
}

func buzzwordsHandler(w http.ResponseWriter, r *http.Request) {
	buzzwordCounter.Inc()
	responseWithText(w, buzzwords.BuzzWords(), 200)
}

func suffixHandler(w http.ResponseWriter, r *http.Request) {
	suffixCounter.Inc()
	responseWithText(w, buzzwords.WithSuffix(), 200)
}

func verbsHandler(w http.ResponseWriter, r *http.Request) {
	verbsCounter.Inc()
	responseWithText(w, buzzwords.WithVerb(), 200)
}

func verbsAndSuffixHandler(w http.ResponseWriter, r *http.Request) {
	verbsAndSuffixCounter.Inc()
	responseWithText(w, buzzwords.WithVerbAndSuffix(), 200)
}
