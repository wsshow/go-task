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

type Task struct {
	name        string
	waitTime    time.Duration
	handlerFunc HandleFunc
	stopSignal  chan bool
}

type TaskManage struct {
	tasks    []Task
	mapTasks map[string]Task
	mu       sync.Mutex
}

func New() *TaskManage {
	return &TaskManage{
		tasks:    make([]Task, 0),
		mapTasks: make(map[string]Task),
		mu:       sync.Mutex{},
	}
}

func (t *TaskManage) Add(f HandleFunc, wt time.Duration) error {
	t.mu.Lock()
	task := Task{
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

func (t *TaskManage) Remove(f HandleFunc) error {
	t.mu.Lock()
	taskName := getFuncName(f)
	task, ok := t.mapTasks[taskName]
	if !ok {
		t.mu.Unlock()
		return errors.New("task name hadn't existed")
	}
	task.stopSignal <- false
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

func (t *TaskManage) Start() {
	for _, task := range t.tasks {
		go func(task Task) {
			for {
				select {
				case <-task.stopSignal:
					log.Println(task.name, "exit")
					return
				case <-time.After(task.waitTime):
					task.handlerFunc()
				}
			}
		}(task)
	}
}

func (t *TaskManage) Stop() {
	for _, task := range t.tasks {
		task.stopSignal <- false
	}
}

func (t *TaskManage) Count() int {
	return len(t.tasks)
}

func getFuncName(f HandleFunc) string {
	return runtime.FuncForPC(reflect.ValueOf(f).Pointer()).Name()
}
