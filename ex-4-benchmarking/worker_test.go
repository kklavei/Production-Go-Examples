package main

import (
	"sync"
	"testing"
)

func Benchmark4Workers(b *testing.B) { benchmarkWorkers(4, b) }

func Benchmark8Workers(b *testing.B) { benchmarkWorkers(8, b) }

func Benchmark16Workers(b *testing.B) { benchmarkWorkers(16, b) }

func Benchmark32Workers(b *testing.B) { benchmarkWorkers(32, b) }

func benchmarkWorkers(numWorkers int, b *testing.B) {
	for n := 0; n < b.N; n++ {
		wg := new(sync.WaitGroup)
		ch := make(chan string)
		go queueMessages(ch)
		runWorkerPool(ch, wg, numWorkers)
	}
}
