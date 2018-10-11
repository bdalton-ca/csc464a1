package main

import (
	"fmt"
	"math/rand"
	//"sync"
	//"sync/atomic"
	"time"
)

type NODE struct {
	data int
	next *NODE
}

const benchmark = false

var list_size = 0
var search_count = 0
var insert_count = 0
var list *NODE
var free_nodes *NODE

var search_mutex chan int
var insert_mutex chan int
var delete_mutex chan int

var no_searchers chan int
var no_inserters chan int

var n_search = 0
var n_insert = 0
var n_delete = 0

const N_NODES = 20
const N_THREADS = 3

func searcher() {
	for {
		if rand.Int()%4096 > 0 {
			continue
		}
		<-search_mutex
		if search_count == 0 {
			<-no_searchers
		}
		search_count++
		search_mutex <- 0

		size := list_size

		if size > 0 {
			r := rand.Int() % size
			node := list
			for i := 0; node.next != nil && i < r; i++ {
				node = node.next
			}
			if !benchmark {
				fmt.Printf("searcher found %d, (size=%d)\n", node.data, size)
			}
		}

		<-search_mutex
		n_search++
		search_count--
		if search_count == 0 {
			no_searchers <- 0
		}
		search_mutex <- 0
	}
}

func inserter() {
	for {
		<-no_inserters
		if N_NODES-list_size > 0 {
			var r int
			if list_size > 0 {
				r = rand.Int() % list_size
			} else {
				r = 0
			}

			node := list
			for i := 0; i < r; i++ {
				node = node.next
			}

			tmp := free_nodes
			free_nodes = free_nodes.next

			if node != nil {
				tmp.next = node.next
				node.next = tmp
			} else {
				tmp.next = nil
				list = tmp
			}

			list_size++

			if !benchmark {
				fmt.Printf("inserter added %d, (size=%d)\n", tmp.data, list_size)
			}

			<-insert_mutex
			n_insert++
			insert_mutex <- 0

		}
		no_inserters <- 0
	}
}
func deleter() {
	for {
		<-no_searchers
		<-no_inserters

		if list_size > 0 {
			var r int
			if list_size-1 > 0 {
				r = rand.Int() % (list_size - 1)
			} else {
				r = 0
			}
			node := list
			for i := 0; i < r; i++ {
				node = node.next
			}

			var tmp *NODE
			if node.next != nil {
				tmp = node.next
				node.next = node.next.next
			} else {
				tmp = node
			}

			tmp.next = free_nodes
			free_nodes = tmp

			if !benchmark {
				fmt.Printf("deleter removed %d\n", tmp.data)
			}

			list_size--

			<-delete_mutex
			n_delete++
			delete_mutex <- 0
		}
		no_searchers <- 0
		no_inserters <- 0
	}
}

func main() {

	fmt.Printf("")

	search_mutex = make(chan int, 1)
	insert_mutex = make(chan int, 1)
	delete_mutex = make(chan int, 1)

	no_searchers = make(chan int)
	no_inserters = make(chan int)

	var nodes [N_NODES]NODE
	for i := 0; i < N_NODES-1; i++ {
		nodes[i].data = i
		nodes[i].next = &nodes[i+1]
	}
	nodes[N_NODES-1].data = N_NODES - 1
	free_nodes = &nodes[0]

	for i := 0; i < N_THREADS*3; i += 3 {
		go searcher()
		go inserter()
		go deleter()
	}

	search_mutex <- 0
	insert_mutex <- 0
	delete_mutex <- 0
	no_searchers <- 0
	no_inserters <- 0

	time.Sleep(time.Millisecond * 2000)

	//fmt.Printf("sid: %d,%d,%d\n", n_search, n_insert, n_delete)

	for {
	}

}
