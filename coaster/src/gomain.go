package main

import (
	"fmt"
	//"math/rand"
	//"sync"
	//"sync/atomic"
	"time"
)

const car_capacity int = 20
const benchmark = false

var n_passengers int = 100
var queue chan int
var unboard chan int

func get_ms() float64 {
	return float64(time.Nanosecond) * float64(time.Now().UnixNano()) / float64(time.Millisecond)
}

func car() {

	queue <- 0
	if !benchmark {
		fmt.Printf("Ready to load\n")
	}
	for i := 0; true; {
		count := <-queue
		if count == car_capacity {
			if !benchmark {
				fmt.Printf("Car is running\n")
				fmt.Printf("Car is unloading\n")
			}
			for count > 0 {
				unboard <- count
				count = <-unboard
			}

			if !benchmark {
				fmt.Printf("Ready to load\n")
			}
			i += car_capacity
			if i == n_passengers {
				break
			}
		}
		queue <- count
	}
}

func passenger() {

	for {
		count := <-queue
		if count < car_capacity {
			count += 1
			if !benchmark {
				fmt.Printf("%d passengers have boarded,\n", count)
			}
			queue <- count

			count = <-unboard
			count -= 1
			if !benchmark {
				fmt.Printf("%d passenger have unboarded\n", car_capacity-count)
			}
			unboard <- count
			break
		} else {
			queue <- count
		}
	}
}

func main() {

	t0 := get_ms()

	queue = make(chan int)
	unboard = make(chan int)
	for i := 0; i < n_passengers; i++ {
		go passenger()
	}
	t1 := get_ms()
	car()
	t2 := get_ms()

	fmt.Printf("init time: %f\n", t1-t0)
	fmt.Printf("proc time: %f\n", t2-t1)
}
