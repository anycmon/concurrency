package main

import "fmt"

func orDone(done <-chan struct{}, in <-chan int) <-chan int {
	out := make(chan int)
	go func() {
		defer close(out)
		for {
			select {
			case <-done:
				return
			case v, ok := <-in:
				if ok == false {
					return
				}

				select {
				case <-done:
					return
				case out <- v:
				}
			}
		}
	}()
	return out
}

func gen(done <-chan struct{}, values ...int) <-chan int {
	out := make(chan int)
	go func() {
		defer close(out)
		for _, v := range values {
			select {
			case <-done:
				close(out)
			case out <- v:
			}
		}
	}()
	return out
}

func multiply(done <-chan struct{}, in <-chan int, m int) <-chan int {
	out := make(chan int)
	go func() {
		defer close(out)
		for v := range orDone(done, in) {
			out <- v * m
		}
	}()
	return out
}
func main() {
	done := make(chan struct{})
	defer close(done)
	for v := range multiply(done, gen(done, 1, 2, 3, 4, 5, 6), 2) {
		fmt.Println(v)
	}
}
