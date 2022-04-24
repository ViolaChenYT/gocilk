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
}

type processor struct {
	ID          int
	currentPrio JobPriority
	p1jobs      *deque.Deque
	p2jobs      *deque.Deque
	p3jobs      *deque.Deque
	p4jobs      *deque.Deque
}

func NewProcessor(id int) Processor {
	p := processor{
		ID:     id,
		p1jobs: deque.New(deqCap, deqMinSize),
		p2jobs: deque.New(deqCap, deqMinSize),
		p3jobs: deque.New(deqCap, deqMinSize),
		p4jobs: deque.New(deqCap, deqMinSize),
	}
	return &p
}

func (p *processor) GetID() int {
	return p.ID
}

func (p *processor) RunJob() error {
	return nil
}

func (p *processor) SetPriority(prio JobPriority) error {
	p.currentPrio = prio
	return nil
}
