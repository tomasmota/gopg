package main

import (
	"fmt"
	"math/rand"
	"time"
)

func main() {
	readCh := receive()
	processCh := process()
	for {
		processCh <- <-readCh
	}
}

func receive() <-chan int {
	// pretend we are listening for something
	ch := make(chan int)
	go func() {
		for {
			time.Sleep(600 * time.Millisecond)
			ch <- rand.Int()
			ch <- rand.Int()
			fmt.Println()
		}
	}()
	return ch
}

func process() chan<- int {
	ch := make(chan int)
	go func() {
		for {
			data := <-ch
			fmt.Println(data)
		}
	}()
	return ch
}
