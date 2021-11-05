package task

import (
	"errors"
	"sync"
	"time"
)

type ITask interface {
	Add(task Task) error
	Remove(taskName string) error
	Start()
	Stop()
	Count() int
}

type Task struct {
	Name        string
	WaitTime    time.Duration
	HandlerFunc func(...interface{})
	StopSignal  chan bool
}

type TaskManage struct {
	Tasks    []Task
	MapTasks map[string]Task
	Mu       sync.Mutex
	Wg       sync.WaitGroup
}

func (t *TaskManage) Add(task Task) error {
	t.Mu.Lock()
	_, ok := t.MapTasks[task.Name]
	if ok {
		t.Mu.Unlock()
		return errors.New("task name had existed")
	}
	t.MapTasks[task.Name] = task
	t.Tasks = append(t.Tasks, task)
	t.Mu.Unlock()
	return nil
}

func (t *TaskManage) Remove(taskName string) error {
	t.Mu.Lock()
	task, ok := t.MapTasks[taskName]
	if !ok {
		t.Mu.Unlock()
		return errors.New("task name hadn't existed")
	}
	task.StopSignal <- false
	delete(t.MapTasks, taskName)
	func() {
		sliceTask := t.Tasks
		cnt := len(sliceTask)
		for i := 0; i < cnt; i++ {
			if sliceTask[i].Name == taskName {
				sliceTask = append(sliceTask[:i], sliceTask[i+1:]...)
				break
			}
		}
		t.Tasks = sliceTask
	}()
	t.Mu.Unlock()
	return nil
}

func (t *TaskManage) Start() {
	for _, task := range t.Tasks {
		go func(task Task) {
			select {
			case <-task.StopSignal:
				return
			case <-time.After(task.WaitTime):
				task.HandlerFunc()
			}
		}(task)
	}
}

func (t *TaskManage) Stop() {
	for _, task := range t.Tasks {
		task.StopSignal <- false
	}
}

func (t *TaskManage) Count() int {
	return len(t.Tasks)
}
