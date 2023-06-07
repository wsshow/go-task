package main

import (
	"log"
	"runtime"
	"time"

	"github.com/wsshow/go-task/task"
)

func main() {
	// 初始化任务池
	var taskM task.ITask = task.New()

	// 向任务池中添加执行函数func1
	err := taskM.Add(func() {
		log.Println("test1...")
	}, time.Second)
	if err != nil {
		log.Panicln(err)
	}
	// 向任务池中添加执行函数func2
	f2 := func() {
		log.Println("test2...")
	}
	err = taskM.Add(f2, time.Second)
	if err != nil {
		log.Panicln(err)
		return
	}

	// 开始异步执行
	taskM.Start()
	time.Sleep(time.Second)

	// 打印当前协程数量
	log.Println("current num of goroutine:", runtime.NumGoroutine())

	// 让任务执行一段时间
	time.Sleep(5 * time.Second)

	// 移除f2函数
	err = taskM.Remove(f2)
	if err != nil {
		log.Panicln(err)
		return
	}

	// 让剩余任务继续执行一段时间
	time.Sleep(5 * time.Second)

	// 打印当前协程数量
	log.Println("current num of goroutine:", runtime.NumGoroutine())

	// 停止所有任务
	taskM.Stop()
	time.Sleep(time.Second)

	// 打印最终协程数量
	log.Println("current num of goroutine:", runtime.NumGoroutine())
}
