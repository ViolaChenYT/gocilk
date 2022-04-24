package main

import (
	"fmt"

	"gocilk/cilk"
)

const (
	duration = 10000
)

func main() {
	server := cilk.NewServer()
	server.Close()
	fmt.Println("fk what just happened")
}
