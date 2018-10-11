5.6 Building H2O

This problem is a basic barrier problem involving two "hydrogen" threads,
and one "oxygen" thread. A thread must block until it can be used
to form a water molecule. This kind of problem could represent a job system
with some sort of dependency constraints.

The implementations given are adaptations of those given in the little book of semaphores.

Correctness
	
	For simplicity's sake I made sure that there would always be an even ratio of
	hydrogrens to oxygens in my algorithm. It is pretty easy to see that the output
	of the program is grouped appropriately int molecules, so the threads must be synchronizing corretly.
	Output files are provided in "cout.txt" and "gout.txt" for C and Go respectively.
	
Comprehensibility

	I tried to use a waitgroup as a barrier, which made Go and C about as one to one as possible,
	but unfortunately Go's waitgroup is not as versatile as pthreads barrier.
	There is no atomic Wait+Done operation in the waitgroup class, and so this makes the 
	problem a more difficult. I replaced the barrier with two channels, one to count the number
	of waiting atoms, and a counter for them to wait on until they had all been assembled.
	
	Aside from that channels and semaphores again map on to eachother pretty well, 
	as in the barbershop problem. I used a mutex from the standard library in Go.
	I could just as easily have used a channel, but I chose not to when I was writing it.
	
Performance

	Atoms		|	C Time (ms)	|	Go Time (ms)
	
	N=10		|	1.0			|	.0
	N=100		|	5.0			|	.0
	N=400		|	18.0		|	4.0
	N=800		|	39.0		|	13.0
	
	Again C is slower than Go, although not quite as drammatically as in the barbershop example.
	