package scheduler

import (
	"log"
	"sync"
	"time"
)

type Scheduler struct {
	task     func()
	taskName string
	trigger  chan struct{}
	active   bool
	mu       *sync.Mutex
}

func NewScheduler(task func(), taskName string) *Scheduler {
	return &Scheduler{
		task:     task,
		taskName: taskName,
		trigger:  make(chan struct{}, 1),
		mu:       &sync.Mutex{},
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
				log.Printf("get request to %s from timer, but already in progress\n", s.taskName)
				s.mu.Unlock()
				continue
			}
			log.Printf("start %s\n", s.taskName)
			s.active = true
			s.mu.Unlock()
			go func(s *Scheduler) {
				defer func(s *Scheduler) {
					log.Printf("end of %s\n", s.taskName)
					s.mu.Lock()
					s.active = false
					s.mu.Unlock()
				}(s)
				s.task()
			}(s)
		case <-s.trigger:
			s.mu.Lock()
			if s.active {
				log.Printf("get request to %s from timer, but already in progress\n", s.taskName)
				s.mu.Unlock()
				continue
			}
			log.Printf("start %s\n", s.taskName)
			s.active = true
			s.mu.Unlock()
			go func(s *Scheduler) {
				defer func(s *Scheduler) {
					log.Printf("end of %s\n", s.taskName)
					s.mu.Lock()
					s.active = false
					s.mu.Unlock()
				}(s)
				s.task()
			}(s)
		}
	}
}

func (s *Scheduler) RunNow() {
	select {
	case s.trigger <- struct{}{}:
	default:
	}
}
