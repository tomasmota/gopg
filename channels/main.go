package main

import (
	"fmt"
	"math/rand"
	"sync"
	"time"
)

func main() {
	ch := make(chan int, 10)
	wg := sync.WaitGroup{}

	f := func(wg *sync.WaitGroup) {
		time.Sleep(time.Second)
		ch <- rand.Int()
		wg.Done()
	}

	for range 5 {
		wg.Add(1)
		go f(&wg)
	}

	wg.Wait()

	for range 5 {
		num := <-ch
		fmt.Println(num)
	}
}

func receiveData() <-chan int {
    // pretend we are listening for something
    ch := make(chan int)
    go func(){
        time.Sleep(time.Second)
        ch <- rand.Int()

    }()
    return ch
}

// func dataProcessor(ch chan<- int) {
//
// 		num := <-ch
// 		fmt.Println(num)
// }
