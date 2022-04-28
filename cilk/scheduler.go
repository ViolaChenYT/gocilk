package cilk

// scheduler of a specific priority level
type Scheduler interface {
	// the scheduler have access to all the deques of all the processors
	NewJob(*Job) error

	// adjust what processors are "mine"
	AdjustProcs() error

	// calculate utilizations from all the processors
	CalculateUtils() float64

	// sum over desires across all the processors under me
	CalculateDesires() float64

	Close() error
}

type scheduler struct {
	// access to all processors
	priority       JobPriority
	processors     map[int]Processor
	myProcessors   map[int]bool
	totalDesire    float32
	receiveJobChan chan *Job
	quit           chan bool
	jobs           map[int]*Job
	// probably some channels or whatever idk
}

func NewScheduler(p JobPriority, procs map[int]Processor) Scheduler {
	sdlr := scheduler{
		priority:       p,
		processors:     procs,
		myProcessors:   make(map[int]bool),
		quit:           make(chan bool),
		receiveJobChan: make(chan *Job, deqCap),
		jobs:           make(map[int]*Job),
	}
	go sdlr.Main()
	return &sdlr
}

func (sdlr *scheduler) Main() error {
	for {
		select {
		case <-sdlr.quit:
			return nil
		case j := <-sdlr.receiveJobChan:
			sdlr.jobs[j.ID] = j
		}
	}
}
func (sdlr *scheduler) NewJob(j *Job) error {
	sdlr.receiveJobChan <- j
	return nil
}
func (sdlr *scheduler) AdjustProcs() error {
	return nil
}

func (sdlr *scheduler) CalculateUtils() float64 {
	//@TO-DO
	return 0.0
}

func (sdlr *scheduler) CalculateDesires() float64 {
	//@TO-DO
	return 0.0
}

func (sdlr *scheduler) Close() error {
	sdlr.quit <- true
	return nil
}
