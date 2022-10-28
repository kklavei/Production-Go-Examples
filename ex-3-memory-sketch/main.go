package main

import (
	"flag"
	"fmt"
	"math/rand"
	"net/http"
	_ "net/http/pprof"
	"os"
	"strings"
	"sync"
	"time"
)

var numWords int

func runWorkerPool(ch chan string, wg *sync.WaitGroup, numWorkers int) {
	wg.Add(numWorkers)
	for i := 0; i < numWorkers; i++ {
		go func() {
			defer wg.Done()
			worker(ch, wg)
		}()
	}
	wg.Wait()
}

// getMessages gets a slice of messages to process.
func getMessages() []string {
	file, _ := os.ReadFile("../datums/melville-moby_dick.txt")
	words := strings.Split(string(file), " ")
	return words
}

// this will block and not close if the len(msgs) is larger than the channel buffer.
func queueMessages(ch chan string) {
	msgs := getMessages()
	for _, msg := range msgs {
		// add messages to string channel
		ch <- msg
	}

	// close the worker channel and signal there won't be any more data
	close(ch)
}

func worker(ch chan string, wg *sync.WaitGroup) {
	var mu sync.Mutex

	for word := range ch {
		if strings.Contains(word, "whal") {
			//fmt.Printf("%s\n", word)
			mu.Lock()
			numWords++
			mu.Unlock()

			// simulate work
			length := time.Duration(rand.Int63n(100))
			time.Sleep(length * time.Millisecond)
		}
	}
}

func main() {
	numWorkers := flag.Int("workers", 1, "number of workers")
	flag.Parse()

	rand.Seed(time.Now().Unix())

	// run pprof
	go func() {
		http.ListenAndServe("localhost:6060", nil)
	}()

	wg := new(sync.WaitGroup)
	ch := make(chan string)
	startTime := time.Now()

	// start the workers in the background and wait for data on the channel
	// we already know the number of workers, we can increase the WaitGroup once
	go queueMessages(ch)
	runWorkerPool(ch, wg, *numWorkers)
	fmt.Printf("Number of words: %d\nTime to process file: %2f seconds\n", numWords, time.Since(startTime).Seconds())
}
