package main

import (
	"fmt"
	"time"
)

func main() {
	ch := make(chan int, 2)
	exit := make(chan struct{})

	go func() {

		for i := 0; i < 3; i++ {
			fmt.Println(time.Now(), i, "sending")
			ch <- i
			fmt.Println(time.Now(), i, "sent")

			time.Sleep(1 * time.Second)
		}

		fmt.Println(time.Now(), "all completed, leaving...")

		close(ch)

	}()

	go func() {
		// Method 1: infinite for with select
		// overcomplicated in this example. It's better for multiple
		// channels cases
		// for {
		// 	select {
		// 	case v, open := <-ch:
		// 		if !open {
		// 			close(exit)
		// 			return
		// 		}

		// 		fmt.Println(time.Now(), "received", v)
		// 	}
		// }

		// Method 2: range statement, better for one channel usage
		for v := range ch {
			fmt.Println(time.Now(), "received", v)
		}

		close(exit)
	}()

	fmt.Println(time.Now(), "waiting for everything to complete...")

	<-exit

	fmt.Println(time.Now(), "exiting")
}
