package main

import (
	"bufio"
	"fmt"
	"log"
	"net"
)

func main() {
	listener, err := net.Listen("tcp", "localhost:8000")
	if err != nil {
		log.Fatal(err)
	}
	go broadcaster()
	for {
		conn, err := listener.Accept()
		if err != nil {
			log.Print(err)
			continue
		}
		go handleConn(conn)
	}
}

type client chan<- string // an outgoing message channel

var (
	entering = make(chan client)
	leaving  = make(chan client)
	messages = make(chan string) // all incoming client messages
)

// 后台运行，等待channel消息到达
func broadcaster() {
	clients := make(map[client]bool) // all connected clients
	for {
		select {
		case msg := <-messages:
			// Broadcast incoming message to all
			// clients' outgoing message channels.
			for cli := range clients {
				cli <- msg
			}
		case cli := <-entering:
			clients[cli] = true

		case cli := <-leaving:
			delete(clients, cli)
			close(cli)
		}
	}
}

func handleConn(conn net.Conn) {
	ch := make(chan string) // outgoing client messages
	go clientWriter(conn, ch)

	who := conn.RemoteAddr().String()
	ch <- "You are " + who
	messages <- who + " has arrived"
	fmt.Println("enter:", ch)

	entering <- ch

	input := bufio.NewScanner(conn)
	for input.Scan() {
		if input.Text() == "exit" {
			ch <- who + " is " + "exit !"
			break
		} else {
			messages <- who + ": " + input.Text()
		}
	}
	// NOTE: ignoring potential errors from input.Err()

	//chSth := client(<-ch) //string can not convert to client ?
	//leaving <- chSth

	leaving <- ch // 当客户端netcat3直接关掉程序时,input.Scan()返回false，由此跳出循环，这里会得到执行
	fmt.Println("leave:", "" == <-ch, ",input.Scan():", input.Scan())
	// 当ch里面没有值，还试图拿值，应该是会阻塞，为何这里ch拿出来的是string的零值？这种情况应该是channel被关闭之后，试图在这里拿值才会产生的零值
	messages <- who + " has left"
	conn.Close()
}

/*
handleConn为每一个客户端创建了一个clientWriter的goroutine，
用来接收向客户端发送消息的channel中的广播消息，并将它们写入到客户端的网络连接。
客户端的读取循环会在broadcaster接收到leaving通知并关闭了channel后终止。
*/
func clientWriter(conn net.Conn, ch <-chan string) {
	for msg := range ch {
		fmt.Fprintln(conn, msg) // NOTE: ignoring network errors
	}
}