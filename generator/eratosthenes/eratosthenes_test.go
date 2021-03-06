package eratosthenes_test

import (
  "testing"
  "github.com/jbuberel/prime_finder/generator/eratosthenes"
  )

func Test100(t *testing.T) {
  r, _ := eratosthenes.GetPrime(int64(100))
  if r != 97 {
    t.Fatalf("Result not equal to 97")
  }
}

func Test1000(t *testing.T) {
  r, _ := eratosthenes.GetPrime(int64(1000))
  if r != 997 {
    t.Fatalf("Result not equal to 997")
  }
}
