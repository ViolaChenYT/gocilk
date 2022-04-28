package main

import (
	"distuv"
	"fmt"
	"time"

	"gocilk/cilk"
)
var id = 1
const (
	lambda   = 2    // 2 jobs per second
	mu       = 3000 // mean job size
)

func main() {
	server := cilk.NewServer()
	ticker := time.NewTicker(5 * time.Second)
	exp := distuv.Exponential{Rate: lambda}
	jobsizes := distuv.Poisson{Lambda: mu}
	timer := time.NewTimer(time.Duration(exp.Rand()) * time.Second)
	// time.Sleep(2 * time.Second)
	
	done := make(chan bool)
	go func() {
		for {
			select {
			case <-done:
				return
			case <-timer.C:
				size := jobsizes.Rand()
				fmt.Println("size", size)
				job := cilk.NewJob(1, id, int(size))
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
