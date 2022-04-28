package cilk

import (
	"time"
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
	SchedulePriorities()
	// get utilization of the scheduler of a particular priority level
	ReportUtilization(JobPriority) float32
	// compute the desires of
	ReportDesire(JobPriority) float32

	//
	Close() error
}

type server struct {
	// denote which level should each processor be working on
	allocation   map[int]JobPriority
	schedulers   map[JobPriority]Scheduler
	utils        map[JobPriority]float64
	desires      map[JobPriority]float64
	processors   map[int]Processor
	time         int
	ticker       *time.Ticker
	quit         chan bool
	incomingJobs chan *Job
	quota        map[JobPriority]float64
}

var priorities = []JobPriority{1, 2, 3, 4}

func NewServer() Server {
	s := server{
		processors:   make(map[int]Processor),
		allocation:   make(map[int]JobPriority),
		schedulers:   make(map[JobPriority]Scheduler),
		utils:        make(map[JobPriority]float64),
		desires:      make(map[JobPriority]float64),
		quota:        make(map[JobPriority]float64),
		incomingJobs: make(chan *Job, deqCap),
		quit:         make(chan bool),
	}
	s.ticker = time.NewTicker(quanta * time.Microsecond)

	for i := 0; i < nProcessor; i++ {
		s.processors[i] = NewProcessor(i)
	}
	for _, i := range priorities {
		s.schedulers[i] = NewScheduler(i, s.processors)
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
		case <-s.ticker.C:
			s.SchedulePriorities()
			continue
		}
	}
}

func (s *server) SchedulePriorities() {
	totalDesire := 0.0
	for _, i := range priorities {
		s.desires[i] = s.schedulers[i].CalculateDesires()
		totalDesire += s.desires[i]
		s.utils[i] = s.schedulers[i].CalculateUtils()
	}
	// if util < rho * quota
	// need to think of how to initialize quota
	//
	// set allocation based on the previous stuff
	//
	// send info to schedulers to adjust processors using AdjustProcs
}

func (s *server) ReportUtilization(p JobPriority) float32 {
	return 0.0
}

func (s *server) ReportDesire(p JobPriority) float32 {
	return 0.0
}

func (s *server) Close() error {
	s.quit <- true
	return nil
}
