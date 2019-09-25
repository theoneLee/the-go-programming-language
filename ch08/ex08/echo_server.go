package main

import (
	"bufio"
	"fmt"
	"log"
	"net"
	"time"
)

func main() {
	listener, err := net.Listen("tcp", "localhost:8000")
	if err != nil {
		log.Fatal(err)
	}

	for {
		conn, err := listener.Accept()
		if err != nil {
			log.Print(err) // e.g., connection aborted
			continue
		}
		handleConnEcho(conn) // handle one connection at a time
	}
}

func handleConnEcho(conn net.Conn) {
	//io.Copy(conn, conn) // NOTE: ignoring errors
	//conn.Close()
	echolist := make(chan string)

	go func() {
		input := bufio.NewScanner(conn)
		for input.Scan() {
			str := input.Text()
			if str == "exit" {
				break
			}
			if len(str) > 0 {
				fmt.Println("input:", str)
				echolist <- str
			}

		}
	}()

	select {
	case <-time.After(10 * time.Second):
		fmt.Println("client time out")
		conn.Close()
	case str := <-echolist:
		fmt.Fprintf(conn, "%v\n", str)
	}

	conn.Close()

}
