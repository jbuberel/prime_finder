// http://en.wikipedia.org/wiki/Sieve_of_Sundaram
package sundaram

import (
	"log"
	"time"
	"math"
)

func printer(primes_channel chan int64, done_channel chan bool) int64 {
	var largest int64
	largest = 0
	for {
		select {
		case p := <-primes_channel:
			if p > largest {
				largest = p
			}
		case <-done_channel:
			log.Printf("Received done signal, largest is %v\n", largest)
			return largest
		}

	}
}

func sieve(primes []bool, prime_channel chan int64, done_channel chan bool) {

	n := len(primes)
	boundi := int(math.Sqrt(float64(n + 1)))
	for i := 1; i <= boundi; i++ {
		boundj := (n - i) / (1 + 2*i)
		for j := i; j <= boundj; j++ {
			remidx := i + j + 2*i*j
			if remidx < n {
				primes[remidx] = false
			}
		}
	}
	for i, v := range primes {
		if v {
			p := int64(2*i + 1)
			if p <= int64(n) {
				prime_channel <- p
			}
		}
	}

	log.Printf("Sending done signal\n")
	done_channel <- true
}


func GetPrime(limit int64) (int64, float64) {

	log.Printf("staring the sieve for largest prime below %v\n", limit)
	prime_channel := make(chan int64, 0)
	done_channel := make(chan bool, 1)
	primes := make([]bool, limit)
	for i, _ := range primes {
		primes[i] = true
	}

	start := float64(time.Now().UnixNano())
	go sieve(primes, prime_channel, done_channel)
	prime := printer(prime_channel, done_channel)
	end := float64(time.Now().UnixNano())
	duration := ((end - start) / 1e9)
	log.Printf("Duration: %fs\n", duration)
	return prime, duration
}
