package main

import (
	"fmt"
	//"math/rand"
	//"sync"
	//"sync/atomic"
	"time"
)

const benchmark = false

const party_max = 5
const n_customers = 400

var queue_size int = 0
var party_size int = 0
var is_party int = 0
var mutex chan int
var block chan int
var wait chan int

func get_ms() float64 {
	return float64(time.Nanosecond) * float64(time.Now().UnixNano()) / float64(time.Millisecond)
}

func customer() {

	<-mutex
	if is_party == 1 {
		queue_size++
		mutex <- 0
		<-block
		queue_size--
	}

	party_size++
	if !benchmark {
		fmt.Printf("customer joins table (party_size=%d)\n", party_size)
	}

	if party_size == party_max {
		is_party = 1
		if !benchmark {
			fmt.Printf("There is a full party of customers\n")
		}
	}

	if is_party == 0 && queue_size > 0 {
		block <- 0
	} else {
		mutex <- 0
	}

	if !benchmark {
		fmt.Printf("customer eats sushi\n")
	}

	<-mutex
	party_size--
	if party_size == 0 {
		is_party = 0
		if !benchmark {
			fmt.Printf("party leaves table\n")
		}
	} else if is_party == 0 {
		if !benchmark {
			fmt.Printf("customer leaves table (party_size=%d)\n", party_size)
		}
	}
	if is_party == 0 && queue_size > 0 {
		block <- 0
	} else {
		mutex <- 0
	}

	wait <- 0
}

func main() {

	t0 := get_ms()

	wait = make(chan int)
	mutex = make(chan int)
	block = make(chan int)

	t1 := get_ms()
	go customer()
	mutex <- 0
	for i := 0; i < n_customers-1; i++ {

		go customer()
	}

	for i := 0; i < n_customers-1; i++ {
		<-wait
	}
	<-mutex

	t2 := get_ms()

	fmt.Printf("init time: %f\n", t1-t0)
	fmt.Printf("proc time: %f\n", t2-t1)
	fmt.Printf("full time: %f\n", t2-t0)
}
