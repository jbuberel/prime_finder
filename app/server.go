package main

import (
    "fmt"
    "net/http"
    "net/url"
    "encoding/json"
    "strconv"
    "log"
    "github.com/jbuberel/prime_finder/generator/eratosthenes"
    "github.com/jbuberel/prime_finder/generator/sundaram"
)

type PrimeGeneratorResult struct {
  Sieve string `json:"sieve"`
	Prime int64	 `json:"prime"`
	Limit int64	 `json:"limit"`
	Duration float64	`json:"compute_time_sec"`
}


// Response to all URLs beginning with '/prime'
// Expects a single query parameter, 'limit' with an integer value
//   Ex: http://host.com/prime?limit=200
// Returns a JSON structure:
//   {"prime":1999,"limit":2000,"compute_time":0.081017088}
func primeHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Add("Content-type", "text/json")

 	log.Printf("URL: %v\n", r.URL)
	m,_  := url.ParseQuery(r.URL.RawQuery)
	log.Printf("Query string map: %v\n", m)
	limit, err := strconv.ParseInt(m["limit"][0], 10, 64)
    if  err != nil {
		log.Printf("Error extracting string from limit parameter\n", err)
		w.WriteHeader(500)
		fmt.Fprintf(w, "Uable to parse value from limit parameter\n")
		return
	} else {
		log.Printf("Using limit value: [%v]\n", limit)
	}

  results := make([]PrimeGeneratorResult,0)
  prime, duration := sundaram.GetPrime(limit)
  results = append(results, PrimeGeneratorResult {
    Sieve: "sundaram",
		Prime: prime,
		Limit: limit,
		Duration: duration,
	})

	prime, duration = eratosthenes.GetPrime(limit)
	results= append(results, PrimeGeneratorResult {
    Sieve: "eratosthenes",
		Prime: prime,
		Limit: limit,
		Duration: duration,
	})
  json_output, err := json.Marshal(results)
	if err != nil {
		log.Printf("Error encoding json data!\n", err)
		w.WriteHeader(500)
		fmt.Fprintf(w, "Uable render JSON output.\n")
		return
	} else {
	    fmt.Fprintf(w, string(json_output))
	}

}

func defaultHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Add("Content-type", "text/plain")
	fmt.Fprintf(w, "Try a limit instead.\n")
}

func init() {
    http.HandleFunc("/prime", primeHandler)
    http.HandleFunc("/", defaultHandler)
}
