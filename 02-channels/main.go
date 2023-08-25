package main

import "fmt"

// NOTE: this code causes a deadlock as there's no
// goroutine, outisde of 'main' interacting with the created channel
// EX:
// fatal error: all goroutines are asleep - deadlock!
func main() {
	var ch chan int
	ch = make(chan int)

	ch <- 10

	v := <-ch
	fmt.Println("received", v)
}
