package main

import (
	"fmt"
	//"math/rand"
	//"sync"
	//"sync/atomic"
	"time"
)

const queue_max = 20
const n_customers = 100

var queue_size int = 0
var newest_customer int = 0
var oldest_customer int = 0
var queue [queue_max]chan int
var qmutex chan int
var benchmark = false

func get_ms() float64 {
	return float64(time.Nanosecond) * float64(time.Now().UnixNano()) / float64(time.Millisecond)
}

func customer() {
	for {
		<-qmutex
		if queue_size < queue_max {
			if !benchmark {
				fmt.Printf("customer enters queue (waiting=%d)\n", queue_size+1)
			}

			id := newest_customer
			newest_customer = (newest_customer + 1) % queue_max
			queue_size++

			qmutex <- 0
			<-queue[id]

			queue_size--

			qmutex <- 0
			break
		} else {
			qmutex <- 0
		}
	}
	if !benchmark {
		fmt.Printf("customer leaves shop\n")
	}
}

func barber() {
	qmutex <- 0

	if !benchmark {
		fmt.Printf("barber enters shop\n")
	}
	for serviced := 0; serviced < n_customers; {
		<-qmutex
		if queue_size > 0 {
			if !benchmark {
				fmt.Printf("barber services customer (waiting=%d)\n", queue_size-1)
			}

			id := oldest_customer
			oldest_customer = (oldest_customer + 1) % queue_max

			serviced++
			queue[id] <- 0
			// customer unlocks queue on wake

		} else {
			qmutex <- 0
		}
	}
	if !benchmark {
		fmt.Printf("barber exits shop\n")
	}
}

func main() {

	t0 := get_ms()
	for i := 0; i < queue_max; i++ {
		queue[i] = make(chan int)
	}
	qmutex = make(chan int)

	for i := 0; i < n_customers; i++ {
		go customer()
	}

	barber()
	t1 := get_ms()

	fmt.Printf("duration: %f", t1-t0)
}
