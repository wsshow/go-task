package task

import (
	"errors"
	"log"
	"reflect"
	"runtime"
	"sync"
	"time"
)

type HandleFunc func()

type ITask interface {
	Add(HandleFunc, time.Duration) error
	Remove(HandleFunc) error
	Start()
	Stop()
	Count() int
}

type task struct {
	name        string
	waitTime    time.Duration
	handlerFunc HandleFunc
	stopSignal  chan bool
}

type taskManage struct {
	tasks    []task
	mapTasks map[string]task
	mu       sync.Mutex
}

func New() *taskManage {
	return &taskManage{
		tasks:    make([]task, 0),
		mapTasks: make(map[string]task),
		mu:       sync.Mutex{},
	}
}

func (t *taskManage) Add(f HandleFunc, wt time.Duration) error {
	t.mu.Lock()
	task := task{
		name:        getFuncName(f),
		waitTime:    wt,
		handlerFunc: f,
		stopSignal:  make(chan bool),
	}
	_, ok := t.mapTasks[task.name]
	if ok {
		t.mu.Unlock()
		return errors.New("task name had existed")
	}
	t.mapTasks[task.name] = task
	t.tasks = append(t.tasks, task)
	t.mu.Unlock()
	return nil
}

func (t *taskManage) Remove(f HandleFunc) error {
	t.mu.Lock()
	taskName := getFuncName(f)
	tk, ok := t.mapTasks[taskName]
	if !ok {
		t.mu.Unlock()
		return errors.New("task name hadn't existed")
	}
	tk.stopSignal <- false
	delete(t.mapTasks, taskName)
	func() {
		sliceTask := t.tasks
		cnt := len(sliceTask)
		for i := 0; i < cnt; i++ {
			if sliceTask[i].name == taskName {
				sliceTask = append(sliceTask[:i], sliceTask[i+1:]...)
				break
			}
		}
		t.tasks = sliceTask
	}()
	t.mu.Unlock()
	return nil
}

func (t *taskManage) Start() {
	for _, taskT := range t.tasks {
		go func(task task) {
			for {
				select {
				case <-task.stopSignal:
					log.Println(task.name, "exit")
					return
				case <-time.After(task.waitTime):
					task.handlerFunc()
				}
			}
		}(taskT)
	}
}

func (t *taskManage) Stop() {
	for _, task := range t.tasks {
		task.stopSignal <- true
	}
}

func (t *taskManage) Count() int {
	return len(t.tasks)
}

func getFuncName(f HandleFunc) string {
	return runtime.FuncForPC(reflect.ValueOf(f).Pointer()).Name()
}
