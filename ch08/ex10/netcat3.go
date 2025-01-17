package main

import (
	"io"
	"log"
	"net"
	"os"
)

func main() {
	conn, err := net.Dial("tcp", "localhost:8000")
	if err != nil {
		log.Fatal(err)
	}
	done := make(chan struct{})
	go func() { //后台运行线程，用来打印从服务端发过来的数据到客户端自己的标准输出上
		//Copy(dst Writer, src Reader) (written int64, err error)
		// 将src的数据拷贝到dst，直到在src上到达EOF或发生错误。返回拷贝的字节数和遇到的第一个错误。
		io.Copy(os.Stdout, conn) // NOTE: ignoring errors
		log.Println("done")
		done <- struct{}{} // signal the main goroutine 服务器链接中断，通知主线程
	}()
	mustCopy2(conn, os.Stdin)
	conn.Close()
	<-done // wait for background goroutine to finish 在这里会阻塞，等待done有值才继续执行，
}

func mustCopy2(dst io.Writer, src io.Reader) {
	if _, err := io.Copy(dst, src); err != nil {
		log.Fatal(err)
	}
}
