package main

import (
	"distuv"
	"fmt"
	"time"

	"gocilk/cilk"
)

const (
	duration = 10000
	lambda   = 2    // 2 jobs per second
	mu       = 3000 // mean job size
)

func main() {
	server := cilk.NewServer()
	ticker := time.NewTicker(5 * time.Second)
	exp := distuv.Exponential{Rate: lambda}
	jobsizes := distuv.Exponential{Rate: 1 / mu}
	timer := time.NewTimer(time.Duration(exp.Rand()) * time.Second)
	// time.Sleep(2 * time.Second)
	id := 1
	done := make(chan bool)
	go func() {
		for {
			select {
			case <-done:
				
			case <-timer.C:
				job := cilk.NewJob(1, id, int(jobsizes.Rand()))
				id++
				server.JobCame(job)
				timer.Reset(time.Duration(exp.Rand()) * time.Second)
			}
		}
	}()
	<- ticker.C
	done <- true
	ticker.Stop()
	fmt.Println("whyyyy")
	server.Close()
}
