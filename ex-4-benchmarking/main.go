package main

import (
	"flag"
	"fmt"
	"log"
	"math/rand"
	"net/http"
	_ "net/http/pprof"
	"os"
	"runtime"
	"runtime/pprof"
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
	cpuprofile := flag.String("cpuprofile", "", "write cpu profile to `file`")
	memprofile := flag.String("memprofile", "", "write memory profile to `file`")
	flag.Parse()

	rand.Seed(time.Now().Unix())

	// run pprof
	if *cpuprofile != "" {
		f, err := os.Create(*cpuprofile)
		if err != nil {
			log.Fatal("could not create CPU profile: ", err)
		}
		defer f.Close() // error handling omitted for example
		if err := pprof.StartCPUProfile(f); err != nil {
			log.Fatal("could not start CPU profile: ", err)
		}
		defer pprof.StopCPUProfile()
	}

	if *memprofile != "" {
		f, err := os.Create(*memprofile)
		if err != nil {
			log.Fatal("could not create memory profile: ", err)
		}
		defer f.Close() // error handling omitted for example
		runtime.GC()    // get up-to-date statistics
		if err := pprof.WriteHeapProfile(f); err != nil {
			log.Fatal("could not write memory profile: ", err)
		}
	}

	go func() {
		log.Println(http.ListenAndServe("localhost:6060", nil))
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
