package main

import (
	"fmt"
	"math/rand"
	"time"
)

func main() {
	ch := make(chan int)
	go receive(ch)
	go process(ch)
	select {}
}

func other() {
}

func receive(ch chan<- int) {
	for {
		time.Sleep(600 * time.Millisecond)
		ch <- rand.Int()
		ch <- rand.Int()
		fmt.Println()
	}
}

func process(ch <-chan int) {
	for {
		data := <-ch
		fmt.Println(data)
	}
}
