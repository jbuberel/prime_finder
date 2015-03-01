// http://en.wikipedia.org/wiki/Sieve_of_Eratosthenes
package eratosthenes

import (
	"log"
	"time"
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
	for i := int64(2); i < int64(len(primes)); i++ {
		if primes[i] {
			prime_channel <- i
			var q int64
			for q = i + i; q < int64(len(primes)); q = q + i {
				primes[q] = false
			}
		}
	}
	log.Printf("Sending done signal\n")
	done_channel <- true
}

func GetPrime(maxval int64) (int64, float64) {

	log.Printf("staring the sieve for largest prime below %v\n", maxval)
	prime_channel := make(chan int64, 0)
	done_channel := make(chan bool, 1)
	primes := make([]bool, maxval)
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
