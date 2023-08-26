package main

import (
	"encoding/csv"
	"fmt"
	"io"
	"os"
)

func main() {
	// 1 - create a channel from the read info in the csv file
	ch1, err := read("file1.csv")
	if err != nil {
		panic(fmt.Errorf("Could not read file1.csv: %v", err))
	}

	// 2 - create the breaking channels
	br1 := breakup("1", ch1)
	br2 := breakup("2", ch1)
	br3 := breakup("3", ch1)

	// 3 - select the breaking channels
	for {
		if br1 == nil && br2 == nil && br3 == nil {
			break
		}

		select {
		case _, open := <-br1:
			if !open {
				br1 = nil
			}
		case _, open := <-br2:
			if !open {
				br2 = nil
			}
		case _, open := <-br3:
			if !open {
				br3 = nil
			}
		}
	}

	fmt.Println("All completed, exiting...")
}

// breakup creates a goroutine that concurrently tries to read from
// channel 'ch' and closes a 'closing channek' at the process end.
func breakup(worker string, ch <-chan []string) chan struct{} {
	// 1 - declares exiting channel
	chE := make(chan struct{})

	// 2 - concurrently tries to read from the provided channel
	go func() {
		for v := range ch {
			fmt.Println(worker, v)
		}
		// 3 - closes the exiting channel
		close(chE)
	}()

	return chE
}

func read(file string) (<-chan []string, error) {
	// 1 - try to open the csv file
	f, err := os.Open(file)
	if err != nil {
		return nil, fmt.Errorf("failed opening file : %s", file)
	}

	// 2 - create the returning channel
	ch := make(chan []string)

	// 3 - create the csv Reader
	cr := csv.NewReader(f)

	// 4 - read each file line
	go func() {
		for {
			record, err := cr.Read()
			if err == io.EOF {
				close(ch)
				return
			}
			ch <- record
		}
	}()
	// 5 - return the channel
	return ch, nil
}
