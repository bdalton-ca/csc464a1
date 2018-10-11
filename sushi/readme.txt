7.1 The sushi bar problem

This problem is based on a sushi bar where guests can
come to sit down, eat, and leave.
The sushibar problem is different from other problems 
of the same variety because there is no mediating thread, 
and it has a dynamic condition where a full table indicates
that the customers are a party, and will leaves a group.
If the table fills, the behaviour of the threads must account
for this without some mediator ensuring this happens.
This could model some constraints on acessing a resource,
or something along those lines.

The code provided is an implementation of the solutions in the little book of semaphores.

Correctness
	
	Output files are provided in "cout.txt" and "gout.txt" for C and Go respectively.
	In both examples parties form and behave correctly, although in C it happens much
	more rarely, likely because there is more time in between customers.
	
Comprehensibility

	These solutions are almost identical between C and Go,
	aside from the obvious swapping of semphores and channels,
	the implementations are largely the same. As in other
	examples, the channels require the redundant passing of a
	zero value to work as semaphores, so this is the only real
	difference.
	
Performance

	Customers	|	C Time (ms) |	Go Time (ms)
	
	N=20		|	1.0 		|	0.0
	N=100		|	5.0			|	1.0
	N=400		|	20.0		|	2.0
	N=800		|	38.0		|	5.0
	
	As seen in other examples, Go again outperforms C.

	