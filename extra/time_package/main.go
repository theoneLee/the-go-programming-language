package main

import (
	"fmt"
	"time"
)

/*
如果select在循环内，第二个case将永远选择不到。因为每次select轮询中，
第一个case都因为2秒而先被选中，使得第二个case的评估总是被中断。
进入下一个select轮询后，又会重新开始评估两个case，分别等待2秒和7秒。
*/
func main() {
	for {
		select {
		case <-time.Tick(2 * time.Second):
			fmt.Println("2 second over:", time.Now().Second())
		case <-time.After(7 * time.Second):
			fmt.Println("5 second over, timeover", time.Now().Second())
			return
		}
	}
}

//不正常执行的原因是因为每次select都会重新评估这些表达式。如果把这些表达式放在select外面，则正常：
// the correct example:

//func main() {
//	tick := time.Tick(1 * time.Second)
//	after := time.After(7 * time.Second)
//	fmt.Println("start second:",time.Now().Second())
//	for {
//		select {
//		case <-tick:
//			fmt.Println("1 second over:", time.Now().Second())
//		case <-after:
//			fmt.Println("7 second over:", time.Now().Second())
//			return
//		}
//	}
//}
