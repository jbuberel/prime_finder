//  +build integration

package generator_test

import (
  "testing"
  e "github.com/jbuberel/prime_finder/eratosthenes"
  s "github.com/jbuberel/prime_finder/sundaram"
  )

func Test100(t *testing.T) {
  eprime, _ := e.GetPrime(100)
  sprime, _ := s.GetPrime(100)
  if eprime != sprime {
    t.Fatalf("Two primes for limit 100 were not equal - %v vs %v\n",eprime, sprime)
  }
}

func Test1000(t *testing.T) {
  eprime, _ := e.GetPrime(1000)
  sprime, _ := s.GetPrime(1000)
  if eprime != sprime {
    t.Fatalf("Two primes for limit 100 were not equal - %v vs %v\n",eprime, sprime)
  }
}


func Test1000000(t *testing.T) {
  eprime, _ := e.GetPrime(1000000)
  sprime, _ := s.GetPrime(1000000)
  if eprime != sprime {
    t.Fatalf("Two primes for limit 100 were not equal - %v vs %v\n",eprime, sprime)
  }
}


func Test1000000000(t *testing.T) {
  eprime, _ := e.GetPrime(1000000000)
  sprime, _ := s.GetPrime(1000000000)
  if eprime != sprime {
    t.Fatalf("Two primes for limit 100 were not equal - %v vs %v\n",eprime, sprime)
  }
}
