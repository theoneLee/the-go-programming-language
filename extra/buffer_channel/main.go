package main

import (
	"fmt"
	"math/rand"
	"sync"
	"time"
)

type Task struct {
	ID         int
	JobID      int
	Status     string
	CreateTime time.Time
}

func (t *Task) run() {
	sleep := rand.Intn(1000)
	time.Sleep(time.Duration(sleep) * time.Millisecond)
	t.Status = "Completed"
}

var wg sync.WaitGroup

// worker的数量，即使用多少goroutine执行任务
const workerNum = 3

func main() {
	wg.Add(workerNum)

	// 创建容量为10的buffered channel
	taskQueue := make(chan *Task, 10)
	//taskQueue := make(chan *Task)

	// 激活goroutine，执行任务
	for workID := 0; workID <= workerNum; workID++ {
		go worker(taskQueue, workID)
	}
	// 将待执行任务放进buffered channel，共15个任务
	for i := 1; i <= 15; i++ {
		fmt.Println("add Task Id:", i)
		taskQueue <- &Task{
			ID:         i,
			JobID:      100 + i,
			CreateTime: time.Now(),
		}
	}
	close(taskQueue)
	wg.Wait()
}

// 从buffered channel中读取任务，并执行任务
func worker(in <-chan *Task, workID int) {
	defer wg.Done()
	for v := range in {
		fmt.Printf("Worker%d: recv a request: TaskID:%d, JobID:%d\n", workID, v.ID, v.JobID)
		v.run()
		fmt.Printf("Worker%d: Completed for TaskID:%d, JobID:%d\n", workID, v.ID, v.JobID)
	}
}
