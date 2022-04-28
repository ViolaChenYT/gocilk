package cilk

import (
	"time"
	"fmt"
)

const (
	nProcessor = 20
	rho        = 2
	delta      = 0.9
)

// top level scheduler,
// keep track of utilization and desier of each priority level
type Server interface {
	// send allocation to each priority scheduler
	JobCame(*Job)
	SchedulePriorities()

	//
	Close() error
}

type server struct {
	// denote which level should each processor be working on
	allocation   map[int]JobPriority
	schedulers   map[JobPriority]Scheduler
	desires      map[JobPriority]int
	processors   map[int]*processor
	time         int
	ticker       *time.Ticker
	quit         chan bool
	incomingJobs chan *Job
	quota        map[JobPriority]int
	alljobs      map[int]int // record flow time
}

var priorities = []JobPriority{1, 2, 3, 4}

func NewServer() Server {
	s := server{
		processors:   make(map[int]*processor),
		allocation:   make(map[int]JobPriority), //priority of each proc
		schedulers:   make(map[JobPriority]Scheduler),
		desires:      make(map[JobPriority]int),
		quota:        make(map[JobPriority]int),
		incomingJobs: make(chan *Job, deqCap),
		quit:         make(chan bool),
	}
	s.ticker = time.NewTicker(quanta * time.Microsecond)

	for i := 0; i < nProcessor; i++ {
		s.processors[i] = NewProcessor(i, s.processors)
	}
	for _, i := range priorities {
		s.schedulers[i] = NewScheduler(i, s.processors)
		s.desires[i] = 0
	}
	go s.Main()
	return &s
}

func (s *server) Main() error {
	for {
		select {
		case <-s.quit:
			return nil
		case job := <-s.incomingJobs:
			job.birthTime = time.Now().UnixMicro()
			sdlr := s.schedulers[job.Prio]
			sdlr.NewJob(job)
			fmt.Println("new job sent")
		case <-s.ticker.C:
			s.SchedulePriorities()
			var already_alloced [5]int

			for pi, p := range s.processors {
				for _, i := range priorities {
					if p.currentPrio == i && already_alloced[i] < s.quota[i] {
						already_alloced[i]++
						s.allocation[pi] = i
					} else if already_alloced[i] < s.quota[i] {
						already_alloced[i]++
						s.allocation[pi] = i
					}
				}
			}
			for _, i := range priorities {
				mp := make(map[int]bool)
				for pi, _ := range s.processors {
					mp[pi] = (s.allocation[pi] == i)
				}
				s.schedulers[i].AdjustProcs(mp)
			}
		}
	}
}

func (s *server) JobCame(job *Job) {
	s.incomingJobs <- job
}

func (s *server) SchedulePriorities() {
	for _, i := range priorities {
		s.desires[i] = s.schedulers[i].CalculateDesires()
	}
	// send info to schedulers to adjust processors using AdjustProcs
	procsLeft := nProcessor
	for _, i := range priorities {
		if s.desires[i] <= s.quota[i] {
			s.quota[i] = s.desires[i]
			procsLeft -= s.quota[i]
		} else if s.desires[i] > procsLeft {
			s.quota[i] = procsLeft
		}
	}
	return
}

func (s *server) Close() error {
	s.quit <- true
	for _, sdlr := range s.schedulers {
		sdlr.Close()
	}
	for _, p := range s.processors {
		p.Close()
	}
	return nil
}
