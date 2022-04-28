package cilk

type JobPriority int

type Job struct {
	Prio      JobPriority
	ID        int
	Size      int
	Done      bool
	birthTime int64 // use time.UnixMicro()
	deathTime int64 // time.UnixMicro()
}

func NewJob(prio JobPriority, id int) *Job {
	return &Job{
		ID:   id,
		Prio: prio,
		Done: false,
	}
}

func (job *Job) GetSize() int {
	return job.Size
}
