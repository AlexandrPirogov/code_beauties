package main

import "fmt"

var multiply = func(values []int, multiplier int) []int {
	multipliedValues := make([]int, len(values))
	for i, v := range values {
		multipliedValues[i] = v * multiplier
	}
	return multipliedValues
}

var add = func(values []int, additive int) []int {
	addedValues := make([]int, len(values))
	for i, v := range values {
		addedValues[i] = v + additive
	}

	return addedValues
}

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


func main() {
	done := make(chan interface{})
  defer close(done)
  intStream := generator(done, 1, 2, 3, 4)
  
  pipeline := multiply(done, add(done, multiply(done, intStream, 2), 1), 2) // <--------------final pipeline
  for v := range pipeline {
   fmt.Println(v)
  }
}
