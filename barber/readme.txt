5.3 The FIFO barbershop

The barbershop problem is a classic problem in concurrency
modelled by customers and a barber which must cut their hair.
The customers represent threaded processes which rely on
some service provided by the barber thread. This FIFO variant
adds the restriction that the customers must be serviced in order of arrival.

The implementations given are adaptations of those given in the little book of semaphores.

Correctness
	
	Output files are provided in "cout.txt" and "gout.txt" for C and Go respectively.
	Both implementations produce effectively the same output, and there appears to be no
	violation of rules, such as customers entering the queue when it full.
	This assumes that the queue is implemented correctly, but it is fairly simple,
	and there don't appear to be any errors either way.

	Mingw32 g++ has a maximum number of semaphores that can be used, a bit below 1024.
	Go doesn't appear to have this limit, or its much larger than I am willing to check.
	This could be important when using many threads, although 1024 hopefully is enough for ordinary applications.
	C does not fail gracefully here, and throws an exception, which causes the program to just hang.
	This will need to be kept in mind for performance benchmarking.
	
Comprehensibility

	It's clear that posting and waiting for messages on a channel
	is effectively the same as doing the same with a semaphore, 
	with the added benefit of being able send values over the channel.
	Unfortunately, this isn't necessary for this example, so the channels
	really only added redundancy. (qmutex <- 0)
	
	Unlike semaphores, unbuffered channels can't be posted to without blocking.
	This means we have to post to qmutex at the start of the barber process	in Go,
	where in C we can just initialize the semaphore to one without blocking.
	This is a valuable feature here, and I think it makes it more difficult to
	write clean code in for the Go example.
	
Performance

	Customers	|	Queue Size	|	C Time (ms)	|	Go Time (ms)
	
	N=10		|	20			|	1.0			|	0.0
	N=100		|	20			|	12.0		|	2.0
	N=400		|	20			|	181.0		|	28.0
	N=800		|	20			|	774.0		|	108.0
	
	N=10		|	100			|	1.0			|	0.0
	N=100		|	100			|	5.0			|	2.0
	N=400		|	100			|	112.0		|	18.0
	N=800		|	100			|	633.0		|	88.0
	
	Because of the semaphore limit, we can only look at values for N < ~1024.
	It looks like the semaphore implementation is on the order of 10x slower than the channel version.
	The queue size appears to be a limiting factor in both cases.