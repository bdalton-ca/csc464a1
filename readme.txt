The code in this repository has been tested on Windows 7 
with mingw32 using g++ 6.3.0 and the Go 1.11 compiler.
The problems in this repository are based on the following
example problems in the Little Book of Semaphores:

	5.3 The FIFO barbershop
	5.6 Building H2
	5.8 The roller coaster problem
	6.1 The search-insert-delete problem
	7.1 The sushi bar problem
	
And also a chat bot problem which will be used for my final project.
	
My intent in this report is to compare the usage of Go's channels with
a more classical semaphore-based implmentations of the examples.
In general my findings were that channels can slightly reduce the amount
of code required over a semaphore implementation, but they can usually
be used interchangably with minor modifications.

In my analysis, I kept my implementations as similar as possible,
so I evaluated the common method rather than comparing differences. 
For performance I recorded the amount of time the entire program took 
to complete, or, where it made more sense, the amount of operations 
that could be completed in a certain amount of time. This was not
exactly relevant for all the examples, because in real world applications
these concurrent tasks have much heavier workloads. Still, I think my analysis 
gives an idea of how effectively the langauge primitives handle concurrency.
With comprehensibility, I found that there were a number of cases
where the only thing that had to change was the syntax difference
between Go's channel notation and C's function calls. 
I decided that I would leave out small syntax differences
in my comparison, because I feel they are largely subjective.
In the end, the only significant difference is found in the roller coaster
problem, where Go's channels do show some benefits over semaphores.

Overall Go appears to be more performant with concurrenct tasks,
and appears to be a good choice for writing optimized concurrent software.
On the other hand C is generally sufficient in this domain, and
since the other differences are negligible, I don't believe channels
alone justify me using Go for my final project. Nevertheless, this
assignment was a good opportunity to experiment with a different language
and evaluate its merits.