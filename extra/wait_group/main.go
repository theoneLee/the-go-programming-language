package main

import (
	"fmt"
	"sync"
)

/*
这两个goroutine执行完全一样的函数代码，它们都接收count这个channel的数据，
但可能是goroutine1先接收到channel中的初始值1，也可能是goroutine2先接收到初始值1。
接收到数据后输出值，并在输出后对数据加1，然后将加1后的数据再次send到channel，
每次send都会将自己这个goroutine阻塞(因为unbuffered channel)，此时另一个goroutine因为等待recv而执行。
当加1后发送给channel的数据为10之后，某goroutine将关闭count channel，该goroutine将退出，wg的计数器减1，
另一个goroutine因等待recv而阻塞的状态将因为channel的关闭而失败，
ok状态码将让该goroutine退出，于是wg的计数器减为0，main goroutine因为wg.Wait()而继续执行后面的代码。
*/

// wg用于等待程序执行完成
var wg sync.WaitGroup

func main() {
	count := make(chan int)

	// 增加两个待等待的goroutines
	wg.Add(2)
	fmt.Println("Start Goroutines")

	// 激活一个goroutine，label："Goroutine-1"
	go printCounts("Goroutine-1", count)
	// 激活另一个goroutine，label："Goroutine-2"
	go printCounts("Goroutine-2", count)

	fmt.Println("Communication of channel begins")
	// 向channel中发送初始数据
	count <- 1

	// 等待goroutines都执行完成
	fmt.Println("Waiting To Finish")
	wg.Wait()
	fmt.Println("\nTerminating the Program")
}
func printCounts(label string, count chan int) {
	// goroutine执行完成时，wg的计数器减1
	defer wg.Done()
	for {
		// 从channel中接收数据
		// 如果无数据可recv，则goroutine阻塞在此
		val, ok := <-count
		if !ok {
			fmt.Println("Channel was closed:", label)
			return
		}
		fmt.Printf("Count: %d received from %s \n", val, label)
		if val == 10 {
			fmt.Printf("Channel Closed from %s \n", label)
			// Close the channel
			close(count)
			return
		}
		// 输出接收到的数据后，加1，并重新将其send到channel中
		val++
		count <- val
	}
}
