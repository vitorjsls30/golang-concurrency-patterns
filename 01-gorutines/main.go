package main

import "fmt"

func main() {
	// NOTE: In this code there's no way to tell main
	// to wait for the hello to finish it's execution
	go hello()
}

func hello() {
	fmt.Println("Yo! Most probably you won't see this message!!")
}
