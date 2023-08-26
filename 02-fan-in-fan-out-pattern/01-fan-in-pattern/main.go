package main

import (
	"encoding/csv"
	"fmt"
	"io"
	"os"
	"sync"
)

// fan-in pattern is used to merge
// multiple chain entries into a single result chain
func main() {
	// 1 - declaring the multiple entry points, the channels themselves...
	ch1, err := read("file1.csv")
	if err != nil {
		panic(fmt.Errorf("Could not read file1: %v", err))
	}

	ch2, err := read("file2.csv")
	if err != nil {
		panic(fmt.Errorf("Could not read file2: %v", err))
	}

	// 2 - creating our clean exit channel
	exit := make(chan struct{})

	// 3 - merging the channels
	// ex 1: using merge1 method...
	// mergedChannels := merge1(ch1, ch2)

	// ex 2: using merge2 method...
	mergedChannels := merge2(ch1, ch2)

	go func() {
		// 4 - reading from the merged channels
		for v := range mergedChannels {
			fmt.Println(v)
		}
		close(exit)
	}()

	// 5 - waiting for the exit channel to be closed
	<-exit

	fmt.Println("All completed, exiting...")
}

// merge1 merges as many channels are provided using a waitGroup for
// the operation synchronization
func merge1(csvChannels ...<-chan []string) <-chan []string {
	// 1 - the waitGroup definition
	var wg sync.WaitGroup

	// 2 - the resulting channel declaration
	out := make(chan []string)

	// 3 - the sender function, extracts the received
	// channel content into the 'out' resulting channel
	send := func(c <-chan []string) {
		for record := range c {
			out <- record
		}
		wg.Done()
	}

	// 4 - add the waitGroup counters based on the received csvChannels argument
	wg.Add(len(csvChannels))

	// 5 - send the current csvChannel data into the out channel
	for _, c := range csvChannels {
		go send(c)
	}

	// 6 - this goroutine keeps waiting for the waitGroup process end
	// so that it can gracefully close the channel usage
	go func() {
		wg.Wait()
		close(out)
	}()

	return out
}

// merge2 merges as many channels are provided using a chan struct{}
// and a counter variable as a mechanism to replace the previous waitGroup usage
func merge2(csvChannels ...<-chan []string) <-chan []string {
	chans := len(csvChannels)
	wait := make(chan struct{}, chans)

	out := make(chan []string)

	send := func(c <-chan []string) {
		defer func() { wait <- struct{}{} }()
		for record := range c {
			out <- record
		}
	}

	for _, c := range csvChannels {
		go send(c)
	}

	go func() {
		for range wait {
			chans--
			if chans == 0 {
				break
			}
		}
		close(out)
	}()

	return out
}

// read reads a given file content records into a result channel
func read(file string) (<-chan []string, error) {
	// 1 - file opening
	f, err := os.Open(file)
	if err != nil {
		return nil, fmt.Errorf("Error opening file: %v", err)
	}

	// 2 - creating the channel to be retuned
	ch := make(chan []string)

	// 3 - creating the content reader
	cr := csv.NewReader(f)

	// 4 - declaring the goroutine that will
	// read the file contents...
	go func() {
		for {
			record, err := cr.Read()
			if err == io.EOF {
				// if end of file, close the channel and return...
				close(ch)
				return
			}
			ch <- record
		}
	}()

	return ch, nil
}
