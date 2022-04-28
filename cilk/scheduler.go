package cilk

import (
	"fmt"
)

// scheduler of a specific priority level
type Scheduler interface {
	// the scheduler have access to all the deques of all the processors
	NewJob(*Job) error

	// adjust what processors are "mine"
	AdjustProcs(map[int]bool) error

	// calculate utilizations from all the processors
	CalculateUtils() int64

	// sum over desires across all the processors under me
	CalculateDesires() int

	Close() error
}

type scheduler struct {
	// access to all processors
	priority       JobPriority
	processors     map[int]*processor
	myProcessors   map[int]bool
	desire         int
	receiveJobChan chan *Job
	quit           chan bool
	jobs           map[int]*Job
	// probably some channels or whatever idk
}

func NewScheduler(p JobPriority, procs map[int]*processor) Scheduler {
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
			fmt.Println("scheduler got job ",j.ID)
		}
	}
}
func (sdlr *scheduler) NewJob(j *Job) error {
	sdlr.receiveJobChan <- j
	return nil
}
func (sdlr *scheduler) AdjustProcs(myprocs map[int]bool) error {
	sdlr.myProcessors = make(map[int]bool)
	for i, b := range myprocs {
		if b {
			sdlr.myProcessors[i] = b
		}
	}
	return nil
}

func (sdlr *scheduler) CalculateUtils() int64 {
	total := 0
	for i, j := range sdlr.jobs {
		if j.Done {
			total += j.Size
			delete(sdlr.jobs, i)
		}
	}
	return int64(total)
}

func (sdlr *scheduler) CalculateDesires() int {
	//@TO-DO
	if sdlr.desire == 0 {
		sdlr.desire = 1
	} else if sdlr.CalculateUtils() < delta*quanta {
		sdlr.desire = sdlr.desire / rho
	} else if len(sdlr.myProcessors) == sdlr.desire {
		sdlr.desire = rho * sdlr.desire
	}
	return sdlr.desire
}

func (sdlr *scheduler) Close() error {
	sdlr.quit <- true
	return nil
}
