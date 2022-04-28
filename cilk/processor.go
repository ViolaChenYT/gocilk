package cilk

import (
	"sync"
	"time"
	"fmt"

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
	deques 			map[JobPriority]*deque.Deque
	processors  map[int]*processor
	quit        chan bool
	currentDeq  *deque.Deque
	mux         sync.Mutex
	newQuanta   chan bool
	busyTime int
}

func NewProcessor(id int, processors map[int]*processor) *processor {
	p := processor{
		ID:         id,
		deques: make(map[JobPriority]*deque.Deque),
		processors: processors,
		quit:       make(chan bool),
		newQuanta:  make(chan bool),
	}
	for _,i := range priorities{
		p.deques[i] = deque.New(deqCap, deqMinSize)
	}
	p.currentPrio = 1
	p.currentDeq = p.deques[p.currentPrio]
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
			if p.currentPrio == 0{ // idle
				continue
			}
			p.currentDeq = p.deques[p.currentPrio]
			p.StealWork()
			p.RunJob()
		case <-p.quit:
			return nil
		}
	}
}

func (p *processor) gotJob(job *Job){
	prio := job.Prio
	p.deques[prio].PushBack(job)
	fmt.Printf("processor %d got job %d of size %d\n", p.ID, job.ID, job.Size)
}

func (p *processor) StealWork(){
	if p.currentDeq.Len() == 0{
		for _, peer := range p.processors{
			peer.mux.Lock()
			if peer.deques[p.currentPrio].Len() > 0{
				job := peer.deques[p.currentPrio].PopBack()
				p.currentDeq.PushFront(job)
				peer.mux.Unlock()
				return
			}
			peer.mux.Unlock()
		}
	}
}

func (p *processor) HandleFork(job *Job) {
	var wg sync.WaitGroup
	for _, child := range job.children {
		wg.Add(1)
		p.currentDeq.PushFront(child)
	}
	wg.Wait()
	// could do more things here
	job.Done = true
}

func (p *processor) NewQ() {
	// fmt.Println("new q")
	p.newQuanta <- true
}

func (p *processor) GetID() int {
	return p.ID
}

func (p *processor) RunJob() error {
	// just reduce the size of job or something?
	timePassed := 0
	for p.currentDeq.Len() != 0 && timePassed <= quanta {
		j := (p.currentDeq.PopFront()).(*Job)
		fmt.Printf("working on job %d\n",j.ID)
		if j.Size < quanta {
			timePassed += j.Size
			j.deathTime = time.Now().UnixMicro()
			j.Done = true
			fmt.Printf("job %d done\n",j.ID)
		} else {
			j.Size = j.Size - quanta
			p.currentDeq.PushFront(j)
		}
	}
	p.busyTime = timePassed
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
