package main

import (
	"fmt"
	"math"
	"sync"
)

// isPrime checks if a number is a prime
func isPrime(n int) bool {
	if n <= 1 {
		return false
	}
	for i := 2; i <= int(math.Sqrt(float64(n))); i++ {
		if n%i == 0 {
			return false
		}
	}
	return true
}

// primeWorker is a worker goroutine that generates prime numbers in the given range
func primeWorker(start, end int, wg *sync.WaitGroup, primeChan chan int) {
	defer wg.Done()
	for i := start; i <= end; i++ {
		if isPrime(i) {
			primeChan <- i
		}
	}
}

// generatePrimes concurrently generates prime numbers up to a given limit
func generatePrimes(limit int, numWorkers int) []int {
	primeChan := make(chan int)
	var wg sync.WaitGroup
	primes := []int{}

	// Divide the range among the workers
	rangeSize := limit / numWorkers
	for i := 0; i < numWorkers; i++ {
		start := i * rangeSize
		end := start + rangeSize - 1
		if i == numWorkers-1 {
			end = limit
		}
		wg.Add(1)
		go primeWorker(start, end, &wg, primeChan)
	}

	// Collect primes in a separate goroutine
	go func() {
		wg.Wait()
		close(primeChan)
	}()

	// Read primes from the channel
	for prime := range primeChan {
		primes = append(primes, prime)
	}

	return primes
}

func main() {
	limit := 1000000
	numWorkers := 10

	primes := generatePrimes(limit, numWorkers)
	fmt.Printf("Generated %d prime numbers up to %d\n", len(primes), limit)
	//fmt.Println(primes)
}
