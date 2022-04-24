package cilk

import (
	"time"
)

const (
	nProcessor = 20
)

// top level scheduler,
// keep track of utilization and desier of each priority level
type Server interface {
	// send allocation to each priority scheduler
	SchedulePriorities()
	// get utilization of the scheduler of a particular priority level
	GetUtilization(Scheduler) float32
	// compute the desires of
	ComputeDesire(JobPriority) float32

	//
	Close() error
}

type server struct {
	// denote which level should each processor be working on
	allocation map[int]JobPriority
	schedulers map[JobPriority]Scheduler
	utils      map[JobPriority]float32
	desires    map[JobPriority]float32
	processors map[int]Processor
	time       int
	ticker     *time.Ticker
	quit       chan bool
}

func NewServer() Server {
	s := server{
		processors: make(map[int]Processor),
		allocation: make(map[int]JobPriority),
		schedulers: make(map[JobPriority]Scheduler),
		utils:      make(map[JobPriority]float32),
		desires:    make(map[JobPriority]float32),
		quit:       make(chan bool),
	}
	s.ticker = time.NewTicker(quanta * time.Microsecond)
	go s.Main()
	return &s
}

func (s *server) Main() error {
	for {
		select {
		case <-s.quit:
			return nil
		case <-s.ticker.C:
			// do stuff here
			continue
		}
	}
}

func (s *server) SchedulePriorities() {

}

func (s *server) GetUtilization(sdlr Scheduler) float32 {
	return 0.0
}

func (s *server) ComputeDesire(p JobPriority) float32 {
	return 0.0
}

func (s *server) Close() error {
	s.quit <- true
	return nil
}
