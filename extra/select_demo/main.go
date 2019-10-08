package main

import (
	"fmt"
	"time"
)

func main() {
	ch1 := make(chan int)
	ch2 := make(chan int)
	go pump1(ch1)
	go pump2(ch2)
	go suck(ch1, ch2)
	time.Sleep(1e9) // main goroutine sleep 1s and exit program
}
func pump1(ch chan int) {
	for i := 0; i <= 30; i++ {
		if i%2 == 0 {
			ch <- i
		}
	}
}
func pump2(ch chan int) {
	for i := 0; i <= 30; i++ {
		if i%2 == 1 {
			ch <- i
		}
	}
}
func suck(ch1 chan int, ch2 chan int) {
	for {
		select {
		case v := <-ch1:
			fmt.Printf("Recv on ch1: %d\n", v)
		case v := <-ch2:
			fmt.Printf("Recv on ch2: %d\n", v)
		}
	}
}
