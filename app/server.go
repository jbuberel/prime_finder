package main

import (
    "fmt"
    "net/http"
    "net/url"
    "encoding/json"
    "strconv"
    "log"
    "log/syslog"
    "os"

    "database/sql"
    _ "github.com/go-sql-driver/mysql"

    "github.com/jbuberel/prime_finder/generator/eratosthenes"
    "github.com/jbuberel/prime_finder/generator/sundaram"

    "github.com/garyburd/redigo/redis"
)

type PrimeGeneratorResult struct {
	Sieve    string  `json:"sieve"`
	Prime    int64   `json:"prime"`
	Limit    int64   `json:"limit"`
	Duration float64 `json:"compute_time_sec"`
  RedisCached bool `json:"redis_cached"`
}

// Response to all URLs beginning with '/prime'
// Expects a single query parameter, 'limit' with an integer value
//   Ex: http://host.com/prime?limit=200
// Returns a JSON structure:
//   {"prime":1999,"limit":2000,"compute_time":0.081017088}
func primeHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Add("Content-type", "text/json")

	log.Printf("URL: %v\n", r.URL)
	m, _ := url.ParseQuery(r.URL.RawQuery)
	log.Printf("Query string map: %v\n", m)
	limit, err := strconv.ParseInt(m["limit"][0], 10, 64)
	if err != nil {
		log.Printf("Error extracting string from limit parameter\n", err)
		w.WriteHeader(500)
		fmt.Fprintf(w, "Uable to parse value from limit parameter\n")
		return
	} else {
		log.Printf("Using limit value: [%v]\n", limit)
	}

  log.Printf("About to retrieve key [%T] with value %v\n", strconv.Itoa(int(limit)), strconv.Itoa(int(limit)))
  redisReturnVal, err := redis.String(redisConn.Do("GET", strconv.Itoa(int(limit))))
  if err != nil {
    log.Printf("redis get error: %v\n", err)
  } else {
    log.Printf("redisReturnVal is [%v] of type [%T]\n", redisReturnVal, redisReturnVal)
  }

  var redisCached bool = false
  if redisReturnVal != "" {
    redisCached = true
  }



	results := make([]PrimeGeneratorResult, 0)
	prime, duration := sundaram.GetPrime(limit)
	results = append(results, PrimeGeneratorResult{
		Sieve:    "sundaram",
		Prime:    prime,
		Limit:    limit,
		Duration: duration,
    RedisCached: redisCached,
	}) 

	prime, duration = eratosthenes.GetPrime(limit)
	results = append(results, PrimeGeneratorResult{
		Sieve:    "eratosthenes",
		Prime:    prime,
		Limit:    limit,
		Duration: duration,
    RedisCached: redisCached,
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

  re, err := redisConn.Do("SET", strconv.Itoa(int(results[0].Limit)), strconv.Itoa(int(results[0].Prime)))
  if err != nil {
    log.Println("redis set error: %v\n", err)
  } else {
    log.Printf("set redis k/v as [%v]/[%v]\n", strconv.Itoa(int(results[0].Limit)), strconv.Itoa(int(results[0].Prime)))
    log.Printf("redis returned: %T with value of %v\n", re, re)
  }


  if err := db.Ping(); err != nil {
    return
  }
  _, err = db.Exec("INSERT INTO primes (limit_val, prime) VALUES (?,?)", results[0].Limit, results[0].Prime)
  if err != nil {
          log.Fatal(err)
  }

}

func defaultHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Add("Content-type", "text/plain")
	fmt.Fprintf(w, "Try a limit instead.\n")
}

func healthHandler(w http.ResponseWriter, r *http.Request) {

  fmt.Fprintf(w, "Database address: %v\n", mysql)
  if err := db.Ping(); err != nil {
    w.Header().Add("Content-type", "text/plain")
  	w.WriteHeader(200)
  	fmt.Fprintf(w, "OK, but no DB Connection\n")

  } else {
    w.Header().Add("Content-type", "text/plain")
  	w.WriteHeader(200)
  	fmt.Fprintf(w, "OK\n")

  }
}

func resultsHandler(w http.ResponseWriter, r *http.Request) {

  w.Header().Add("Content-type", "text/plain")
  w.WriteHeader(200)




  if err := db.Ping(); err != nil {
    fmt.Fprintf(w,"No DB connection\n")
    return
  }
  rows, err := db.Query("SELECT limit_val, prime FROM primes")
  if err != nil {
          log.Println(err)
  }
  defer rows.Close()
  fmt.Fprintf(w, "MySQL Primes are:\n\n")
  for rows.Next() {
    var limit int
    var prime int
    if err := rows.Scan(&limit, &prime); err != nil {
            log.Println(err)
    }
    fmt.Fprintf(w,"Limit: %v - Prime: %v\n", limit, prime)
  }
  if err := rows.Err(); err != nil {
          log.Println(err)
  }


}

var db  *sql.DB
var redisConn redis.Conn

var mysql string

func init() {
  var err error


  mysql = os.Getenv("MYSQL")
  if mysql == "" {
    mysql = "[2001:4860:4864:1:3907:3b3d:5490:9e64]"
  }

  log.Printf("Connecting to database address: %v", mysql)

  db, err = sql.Open("mysql", "service:abc123@tcp("+ mysql +":3306)/primes_schema")
  if err != nil {
    log.Printf("Error connecting: %v", err)
  } else {
    err = db.Ping()
    if err != nil {
      log.Printf("Unable to connect to database: %v\n", err)
    } else {
      log.Printf("Successfully connected to mysql database\n")
    }
  }


  redisConn, err = redis.Dial("tcp", ":6379")
  if err != nil {
    log.Println(err)
  }
  log.Printf("redisConn: %v\n", redisConn)


  sysLogger, err := syslog.New(syslog.LOG_NOTICE, "prime_finder")
  if err == nil {
    log.SetOutput(sysLogger)
  }


}

func main() {
  defer func() {
        if r := recover(); r != nil {
            fmt.Println("Recovered in f", r)
        }
    }()
    defer redisConn.Close()
    defer db.Close()

    http.HandleFunc("/prime", primeHandler)
    http.HandleFunc("/_ah/health", healthHandler)
    http.HandleFunc("/results", resultsHandler)
    http.HandleFunc("/", defaultHandler)

    port := os.Getenv("PORT")
    if port == "" {
      port = "8080"
    }
    log.Printf("Listening on port %v\n", port)
    http.ListenAndServe(":" + port, nil)
}
