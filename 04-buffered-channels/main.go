package main

import (
	"fmt"
	"time"
)

func main() {
	ch := make(chan int, 2)

	go func() {
		for i := 0; i < 3; i++ {
			fmt.Println(time.Now(), i, "sending")
			ch <- i
			fmt.Println(time.Now(), i, "sent")
		}

		// NOTE: there could be cases where this message is not shown
		// solved in future example
		fmt.Println(time.Now(), "all complete")
	}()

	time.Sleep(2 * time.Second)

	fmt.Println("wainting for messages...")

	fmt.Println(time.Now(), "received", <-ch)
	fmt.Println(time.Now(), "received", <-ch)
	fmt.Println(time.Now(), "received", <-ch)

	fmt.Println(time.Now(), "exiting...")
}
