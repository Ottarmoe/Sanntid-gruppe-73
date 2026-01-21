// Use `go run foo.go` to run your program

package main

import (
	. "fmt"
	"runtime"
	//"time"
)

func iterating(servin chan int, servout chan int, c chan int, it int) {
	//TODO: decrement i 1000000 times)
	val := 0
	for j := 0; j < 1000000; j++ {
		val = <-servout
		val += it
		servin <- val
	}
	Println(val)
	c <- 0
}

func intserver(in chan int, out chan int) {
	i := 0
	for {
		out <- i
		i = <-in
	}
}

func main() {
	// What does GOMAXPROCS do? What happens if you set it to 1?
	runtime.GOMAXPROCS(3)
	c := make(chan int)
	servin := make(chan int)
	servout := make(chan int)

	// TODO: Spawn both functions as goroutines
	go intserver(servin, servout)
	go iterating(servin, servout, c, -1)
	go iterating(servin, servout, c, 1)

	<-c
	<-c

	Println("The magic number is:", <-servout)
}
