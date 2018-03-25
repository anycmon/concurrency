package main

import (
	"fmt"
	"time"
)

func main() {
	doWork := func(done <-chan struct{}, pulseInterval time.Duration) (<-chan struct{}, <-chan time.Time) {
		heartbeat := make(chan struct{})
		results := make(chan time.Time)

		go func() {
			//defer close(heartbeat)
			//defer close(results)

			pulse := time.Tick(pulseInterval)
			workGen := time.Tick(2 * pulseInterval)

			sendPulse := func() {
				select {
				case heartbeat <- struct{}{}:
				default: // do not block on heartbeat
				}
			}

			sendResult := func(r time.Time) {
				select {
				case <-done:
					return
				case <-pulse:
					sendPulse()
				case results <- r:
					return
				}
			}

			for i := 0; i < 2; i++ {
				select {
				case <-done:
					return
				case <-pulse:
					sendPulse()
				case r := <-workGen:
					sendResult(r)
				}
			}
		}()
		return heartbeat, results
	}

	done := make(chan struct{})
	time.AfterFunc(10*time.Second, func() { defer close(done) })

	pulseInterval := 2 * time.Second
	heartbeat, results := doWork(done, pulseInterval/2)

	for {
		select {
		case _, ok := <-heartbeat:
			if ok == false {
				return
			}
			fmt.Println("pulse")
		case r, ok := <-results:
			if ok == false {
				return
			}
			fmt.Printf("result %v\n", r.Second())
		case <-time.After(pulseInterval):
			fmt.Println("worker goroutine is not healthy!")
			return
		}
	}
}
