package cilk

type JobPriority int

const (
	level1 JobPriority = iota
	level2
	level3
	level4
)

type Job struct {
	Prio JobPriority
	ID   int
	Size int
	Done bool
}

func NewJob(prio JobPriority, id int) *Job {
	return &Job{
		ID:   id,
		Prio: prio,
		Done: false,
	}
}
