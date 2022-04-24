package cilk

// scheduler of a specific priority level
type Scheduler interface {
	// the scheduler have access to all the deques of all the processors

	// adjust what processors are "mine"
	AdjustProcs() error

	// calculate utilizations from all the processors
	CalculateUtils() float32

	// sum over desires across all the processors under me
	CalculateDesires() float32
}

type scheduler struct {
	// access to all processors
	priority     JobPriority
	processors   map[int]Processor
	myProcessors map[int]bool
	totalDesire  float32
	// probably some channels or whatever idk
}

func NewScheduler(p JobPriority, procs map[int]Processor) Scheduler {
	sdlr := scheduler{
		priority:     p,
		processors:   procs,
		myProcessors: make(map[int]bool),
	}
	return &sdlr
}

func (sdlr *scheduler) AdjustProcs() error {
	return nil
}

func (sdlr *scheduler) CalculateUtils() float32 {
	return 0.0
}

func (sdlr *scheduler) CalculateDesires() float32 {
	return 0.0
}
