package timer

import (
	"sync"
	"time"
)

type TaskFunc func(timerID string, firedAt time.Time, counter int64)

type Timer struct {
	id       string
	interval time.Duration
	mu       sync.RWMutex
	taskMap  map[string]TaskFunc
	stop     chan bool
}

func NewTimer(id string, interval time.Duration) *Timer {
	return &Timer{
		id:       id,
		interval: interval,
		mu:       sync.RWMutex{},
		taskMap:  map[string]TaskFunc{},
		stop:     make(chan bool),
	}
}

func (t *Timer) AddTask(taskId string, taskFunc TaskFunc) {
	t.mu.Lock()
	defer t.mu.Unlock()

	if _, exists := t.taskMap[taskId]; exists {
		return
	}

	t.taskMap[taskId] = taskFunc
}

func (t *Timer) RemoveTask(taskId string) {
	t.mu.Lock()
	defer t.mu.Unlock()

	if _, exists := t.taskMap[taskId]; exists {
		delete(t.taskMap, taskId)
	}
}

func (t *Timer) getTasks() []TaskFunc {
	var tasks []TaskFunc

	t.mu.RLock()
	for _, taskFunc := range t.taskMap {
		tasks = append(tasks, taskFunc)
	}
	t.mu.RUnlock()

	return tasks
}

func (t *Timer) Start() {
	go func() {
		ticker := time.NewTicker(t.interval)
		var counter int64

	TickerLoop:
		for {
			select {
			case <-t.stop:
				break TickerLoop
			case ti := <-ticker.C:
				tasks := t.getTasks()
				for _, taskFunc := range tasks {
					counter = counter + 1
					taskFunc(t.id, ti, counter)
				}
			}
		}

		ticker.Stop()
	}()
}

func (t *Timer) Stop() {
	t.stop <- true
}
