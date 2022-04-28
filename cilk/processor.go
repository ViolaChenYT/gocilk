package cilk

import (
	"github.com/gammazero/deque"
)

const (
	deqCap     = 2048
	deqMinSize = 32
	quanta     = 500
)

type Processor interface {
	GetID() int
	RunJob() error
	SetPriority(JobPriority) error
	Close() error
}

type processor struct {
	ID          int
	currentPrio chan JobPriority
	p1jobs      *deque.Deque
	p2jobs      *deque.Deque
	p3jobs      *deque.Deque
	p4jobs      *deque.Deque
	quit        chan bool
	currentDeq  *deque.Deque
}

func NewProcessor(id int) Processor {
	p := processor{
		ID:          id,
		currentPrio: make(chan JobPriority),
		p1jobs:      deque.New(deqCap, deqMinSize),
		p2jobs:      deque.New(deqCap, deqMinSize),
		p3jobs:      deque.New(deqCap, deqMinSize),
		p4jobs:      deque.New(deqCap, deqMinSize),
		quit:        make(chan bool),
	}
	go p.Main()
	return &p
}

func (p *processor) Main() error {
	for {
		select {
		case prio := <-p.currentPrio:
			switch prio {
			case 1:
				p.currentDeq = p.p1jobs
			case 2:
				p.currentDeq = p.p2jobs
			case 3:
				p.currentDeq = p.p3jobs
			case 4:
				p.currentDeq = p.p4jobs
			}
		case <-p.quit:
			return nil
		}
	}
}

func (p *processor) GetID() int {
	return p.ID
}

func (p *processor) RunJob() error {
	// just reduce the size of job or something?
	return nil
}

func (p *processor) SetPriority(prio JobPriority) error {
	p.currentPrio <- prio
	return nil
}

func (p *processor) Close() error {
	p.quit <- true
	return nil
}
