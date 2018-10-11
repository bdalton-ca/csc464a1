5.8 The roller coaster problem

The coaster problem is based on having several threads synchronize
on a single thread. The car thread allows passengers to board, performs 
its task of travelling around the tracks, and then the passengers can 
unboard and continue as normal. The car waits until it has a full load before departing. 
The problem bears some similarity to a combination of the H2O problem and 
the barbershop problem. Again, it mostly models some sort of job system with dependencies. 

The code provided is my own implementation.

Correctness
	
	Unlike the H2O and barbershop problems I started working on this solution
	using Go's channels to see how they would influence the ultimate design.
	My knowledge of channels is not the best, so it is not at all optimal.
	The main problem with the implementation is that the threads spin and check
	a channel for a value. In hindsight this seems to be a poor choice for performance,
	but Go apparently guarantees a randomized receiver is chosen, so I believe this
	prevents starvation, and isn't incorrect. I believe C pthread semaphores are
	scheduler dependent, so we do not have this guaraentee. This has been fixed
	by a small change to the implementation using two semaphores, which is described
	int comprehensibility.
	
	Output files are provided in "cout.txt" and "gout.txt" for C and Go respectively.
	We can see that the threads produce correct output, and the correct number of
	passengers board and unboard at the right time between rides.
	
Comprehensibility

	Starting the design using channels instead of semaphores has made the code
	a bit more simple and terse in the Go example. In Go we use two channels to
	pass the number of passengers around between threads, whereas we use
	two globals and two mutexes in C. These effectively work the same way,
	though the Go code is much cleaner and there is no confusion about
	how the two globals variables need to be set and reset.
	
	The C code requires an extra two mutexes to prevent the spinning in the Go
	implementation, which causes starvation. The spinning mechanism is replaced
	with two semaphores, where the final passenger signals the car to depart,
	and the car signals and the passengers to unboard. This is a small difference,
	but it does add complexity to implementation.
	
Performance

	Passengers	|	C init (ms)	|	C Time (ms)	|	Go init		|	Go Time (ms)
	
	N=20		|	2.0			|	0.0			|	0.0			|	0.0
	N=100		|	3.0			|	2.0			|	0.0			|	0.0
	N=400		|	20.0		|	15.0		|	1.0			|	2.0
	N=800		|	47.0		|	54.0		|	4.0			|	7.0

	Again Go is about 10x faster than C in this implementation.
	I decided to split my benchmarking between thread initialization time
	and the actual processing time. This reveals a bit more information,
	as C's initialization time is a larger percentage of the overall time 
	than Go. This is not an unexpected result since Go's threads are known
	to be lightweight. 