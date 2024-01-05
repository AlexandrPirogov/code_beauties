package main

import "fmt"

// generating stream-like channel
var generator = func(done <-chan interface{}, integers ...int) <-chan int {
	intStream := make(chan int)
	go func() {
		defer close(intStream)
		for _, i := range integers {
			select {
			case <-done:
				return
			case intStream <- i:
			}
		}
	}()
	return intStream
}

// pipeline stage
var fmap = func(
	done <-chan interface{},
	intStream <-chan int,
	f func(given, change int) int,
	change int,
) <-chan int {
	changedStream := make(chan int)
	go func() {
		defer close(changedStream)
		for i := range intStream {
			select {
			case <-done:
				return
			case changedStream <- f(i, change):
			}
		}
	}()
	return changedStream
}


var mult = func(given, change int) int {
	return given * change
}

var negMult = func(given, change int) int {
	return given * -change
}

var substract = func(given, change int) int {
	return given - change
}

// high-order function
var pipeline = func(stages []func(<-chan int) <-chan int, stream <-chan int) <-chan int {
	prev := stream
	for _, stage := range stages {
		prev = stage(prev)
	}

	return prev
}

func main() {
	done := make(chan interface{})
	defer close(done)

	intStream := generator(done, 1, 2, 3, 4, 5)
	stages := prepareStages(done)
	result := pipeline(stages, intStream)

	for v := range result {
		fmt.Println(v)
	}

}

// Preparing stages
func prepareStages(done <-chan interface{}) []func(<-chan int) <-chan int {
	stage1 := func(stream <-chan int) <-chan int {
		return fmap(done, stream, mult, 2)
	}

	stage2 := func(stream <-chan int) <-chan int {
		return fmap(done, stream, negMult, 1)
	}

	stage3 := func(stream <-chan int) <-chan int {
		return fmap(done, stream, substract, 1)
	}

	stages := []func(<-chan int) <-chan int{stage1, stage2, stage3}
	return stages
}
