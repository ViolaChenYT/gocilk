package cilk

import (
	"fmt"
)

// Second level scheduler
// scheduler of a specific priority level
type Scheduler interface {
	// the scheduler have access to all the deques of all the processors
	NewJob(*Job) error

	// adjust what processors are "mine"
	AdjustProcs(map[int]bool) error

	// calculate utilizations from all the processors under me
	CalculateUtils() int64

	// adjust desire of the priority level
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
			p := sdlr.processors[j.ID % nProcessor]
			p.gotJob(j)
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
		// fmt.Println(i,b)
		if b {
			//fmt.Println(i)
			sdlr.myProcessors[i] = b
			p := sdlr.processors[i]
			p.SetPriority(sdlr.priority)
			p.NewQ()
		}
	}
	// fmt.Println("got ", sdlr.priority, len(sdlr.myProcessors))
	return nil
}

func (sdlr *scheduler) CalculateUtils() int64 {
	total := 0
	for i, j := range sdlr.jobs {
		if j.Done {
			delete(sdlr.jobs, i)
		}
	}
	for i,_ := range sdlr.myProcessors{
		p := sdlr.processors[i]
		if p.busyTime != 0{
			fmt.Println("busy time", p.busyTime)
		}
		total += p.busyTime
	}
	// fmt.Println("utils", total)
	return int64(total)
}

func (sdlr *scheduler) CalculateDesires() int {
	utils := sdlr.CalculateUtils()
	if sdlr.desire == 0 {
		sdlr.desire = 1
	} else if utils < delta * quanta {
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
