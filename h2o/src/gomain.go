package main

import (
	"fmt"
	"math/rand"
	"sync"
	//"sync/atomic"
	"time"
)

func get_ms() float64 {
	return float64(time.Nanosecond) * float64(time.Now().UnixNano()) / float64(time.Millisecond)
}

var wait chan int

var n_oxygen_waiting = 0
var n_hydrogen_waiting = 0
var hydrogen_ch chan int
var oxygen_ch chan int
var mutex sync.Mutex

//var barrier sync.WaitGroup
var barrier chan int
var bwait_ch chan int

const benchmark = false

func oxygen(id int) {

	mutex.Lock()
	n_oxygen_waiting += 1

	if n_hydrogen_waiting >= 2 {

		//barrier.Add(3)
		hydrogen_ch <- 0
		hydrogen_ch <- 0
		oxygen_ch <- 0
		n_hydrogen_waiting -= 2
		n_oxygen_waiting -= 1

	} else {
		mutex.Unlock()
	}

	<-oxygen_ch
	if !benchmark {
		fmt.Printf("o%d\n", id)
	}

	n_waiting := <-bwait_ch
	n_waiting += 1
	if n_waiting == 3 {
		barrier <- 0
		barrier <- 0
		bwait_ch <- 0
	} else {
		bwait_ch <- n_waiting
		<-barrier
	}

	if !benchmark {
		fmt.Printf("H2O formed\n")
	}
	mutex.Unlock()
	wait <- id
}

func hydrogen(id int) {

	mutex.Lock()
	n_hydrogen_waiting += 1

	if n_hydrogen_waiting >= 2 && n_oxygen_waiting >= 1 {

		//barrier.Add(3)
		hydrogen_ch <- 0
		hydrogen_ch <- 0
		oxygen_ch <- 0
		n_hydrogen_waiting -= 2
		n_oxygen_waiting -= 1

	} else {
		mutex.Unlock()
	}

	<-hydrogen_ch
	if !benchmark {
		fmt.Printf("h%d\n", id)
	}

	n_waiting := <-bwait_ch
	n_waiting += 1

	if n_waiting == 3 {
		barrier <- 0
		barrier <- 0
		bwait_ch <- 0
	} else {
		bwait_ch <- n_waiting
		<-barrier
	}

	wait <- id
}

func main() {

	t0 := get_ms()

	//const N_ATOMS = 9
	const N_ATOMS = 99
	//const N_ATOMS = 399
	//const N_ATOMS = 801
	atom_funcs := []func(int){hydrogen, oxygen}
	var atom_types [N_ATOMS]int
	rand.Seed(time.Now().UnixNano())

	mutex = sync.Mutex{}
	//barrier = sync.WaitGroup{}
	barrier = make(chan int)
	wait = make(chan int)           //, N_ATOMS)
	hydrogen_ch = make(chan int, 2) //, 2*N_ATOMS/3)
	oxygen_ch = make(chan int, 1)   //, N_ATOMS/3)
	bwait_ch = make(chan int, 1)
	bwait_ch <- 0

	for i := 0; i < N_ATOMS; i += 3 {
		atom_types[i+0] = 0
		atom_types[i+1] = 0
		atom_types[i+2] = 1
	}

	for i := 0; i < N_ATOMS; i++ {
		for j := i + 1; j < N_ATOMS-1; j++ {
			dst := j + rand.Int()%(N_ATOMS-j)
			tmp := atom_types[i]
			atom_types[i] = atom_types[dst]
			atom_types[dst] = tmp
		}
	}

	for i := 0; i < N_ATOMS; i++ {
		go atom_funcs[atom_types[i]](i)
	}

	for i := 0; i < N_ATOMS; i++ {
		<-wait
	}
	t1 := get_ms()
	fmt.Printf("duration %f", t1-t0)
}
