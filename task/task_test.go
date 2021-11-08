package task

import (
	"log"
	"runtime"
	"strings"
	"testing"
	"time"
)

func Calculate() {
	log.Println("...")
}

var GetFuncName = getFuncName

func TestGetFuncName(t *testing.T) {
	a := GetFuncName(Calculate)
	if !strings.Contains(a, "Calculate") {
		t.Error("Result not match:", a, "Calculate")
	}
	log.Println(GetFuncName(func() {
		log.Println("test1...")
	}))
	log.Println(GetFuncName(func() {
		log.Println("test2...")
	}))
}

func TestTaskManage_workflow(t *testing.T) {
	var taskM ITask = New()
	err := taskM.Add(func() {
		log.Println("test1...")
	}, time.Second)
	if err != nil {
		t.Error(err)
	}

	f2 := func() {
		log.Println("test2...")
	}

	err = taskM.Add(f2, time.Second)
	if err != nil {
		t.Error(err)
		return
	}
	taskM.Start()
	time.Sleep(time.Second)
	log.Println("current num of goroutine:", runtime.NumGoroutine())
	time.Sleep(5 * time.Second)
	err = taskM.Remove(f2)
	if err != nil {
		t.Error(err)
		return
	}
	time.Sleep(5 * time.Second)
	log.Println("current num of goroutine:", runtime.NumGoroutine())
	taskM.Stop()
	time.Sleep(time.Second)
	log.Println("current num of goroutine:", runtime.NumGoroutine())
}
