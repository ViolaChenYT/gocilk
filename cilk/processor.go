package cilk

import (
	"sync"
	"time"

	"github.com/gammazero/deque"
)

const (
	deqCap     = 2048
	deqMinSize = 32
	quanta     = 500
)

type processor struct {
	ID          int
	currentPrio JobPriority
	p1jobs      *deque.Deque
	p2jobs      *deque.Deque
	p3jobs      *deque.Deque
	p4jobs      *deque.Deque
	processors  map[int]*processor
	quit        chan bool
	currentDeq  *deque.Deque
	mux         sync.Mutex
	newQuanta   chan bool
}

func NewProcessor(id int, processors map[int]*processor) *processor {
	p := processor{
		ID:         id,
		p1jobs:     deque.New(deqCap, deqMinSize),
		p2jobs:     deque.New(deqCap, deqMinSize),
		p3jobs:     deque.New(deqCap, deqMinSize),
		p4jobs:     deque.New(deqCap, deqMinSize),
		processors: processors,
		quit:       make(chan bool),
		newQuanta:  make(chan bool),
	}
	p.currentPrio = 1
	p.currentDeq = p.p1jobs
	go p.Main()
	return &p
}

func (p *processor) lock() error {
	p.mux.Lock()
	return nil
}

func (p *processor) unlock() error {
	p.mux.Unlock()
	return nil
}

func (p *processor) Main() error {
	for {
		select {
		case <-p.newQuanta:
			switch p.currentPrio {
			case 1:
				p.currentDeq = p.p1jobs
			case 2:
				p.currentDeq = p.p2jobs
			case 3:
				p.currentDeq = p.p3jobs
			case 4:
				p.currentDeq = p.p4jobs
			}
			// start working on jobs
			p.RunJob()
		case <-p.quit:
			return nil
		}
	}
}

func (p *processor) NewQ() {
	p.newQuanta <- true
}

func (p *processor) GetID() int {
	return p.ID
}

func (p *processor) RunJob() error {
	// just reduce the size of job or something?
	timePassed := 0
	for p.currentDeq.Len() != 0 && timePassed <= quanta {
		j := (p.currentDeq.PopBack()).(*Job)
		if j.Size < quanta {
			timePassed += j.Size
			j.deathTime = time.Now().UnixMicro()
			j.Done = true
		} else {
			j.Size = j.Size - quanta
			p.currentDeq.PushBack(j)
		}
	}
	return nil
}

func (p *processor) SetPriority(prio JobPriority) error {
	p.currentPrio = prio
	return nil
}

func (p *processor) Close() error {
	p.quit <- true
	return nil
}
