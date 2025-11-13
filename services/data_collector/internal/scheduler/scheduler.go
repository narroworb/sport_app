package scheduler

import (
	"log"
	"sync"
	"time"
)

type Scheduler struct {
	task    func()
	trigger chan struct{}
	active  bool
	mu      sync.Mutex
}

func NewScheduler(task func()) *Scheduler {
	return &Scheduler{
		task:    task,
		trigger: make(chan struct{}, 1),
	}
}

func (s *Scheduler) Start(period time.Duration) {
	ticker := time.NewTicker(period)

	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			s.mu.Lock()
			if s.active {
				log.Println("get request to update from timer, but already in progress")
				s.mu.Unlock()
				continue
			}
			log.Println("start update")
			s.active = true
			s.mu.Unlock()
			go func() {
				defer func() { 
					s.mu.Lock()
					s.active = false
					s.mu.Unlock()
				}()
				s.task()
			}()
		case <-s.trigger:
			s.mu.Lock()
			if s.active {
				log.Println("get request to update from api, but already in progress")
				s.mu.Unlock()
				continue
			}
			log.Println("start update")
			s.active = true
			s.mu.Unlock()
			go func() {
				defer func() { 
					s.mu.Lock()
					s.active = false
					s.mu.Unlock()
				}()
				s.task()
			}()
		}
	}
}

func (s *Scheduler) RunNow() {
	select {
	case s.trigger <- struct{}{}:
	default:
	}
}
