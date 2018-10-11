6.1 The search-insert-delete problem

This problem is effectively just creating a thread safe linked list.
It's applications are fairly obvious, and it can be used for
anything a linked list might be used for under normal circumstances.

This implementation is an adaptation of the little book of semaphores solution.

Correctness
	
	The correctness of this method is clear because it is simply
	a data structure with a lock. In the case of the searcher,
	the mutex prevents the delete from deleting the searchers current node,
	which is the only thing that could happen to it. The inserter and deleter
	are made mutually exclusive, so this is essentially making access serial but
	allowing each thread to synchronize with each other.
	
	The main problem with this method is that the searcher threads can starve the other
	threads since it is much lighter weight. To compensate, I've created some busy
	work for the searchers. In reality, it would be unlikely that you would 
	need to constantly search without using the data that you search for.
	
	Output files are provided in "cout.txt" and "gout.txt" for C and Go respectively.
	
Comprehensibility

	These implementations maps one to one with channels and semaphores.
	The only real differences come from non-concurrency language features,
	so I'll won't talk about them here.
	
Performance

	Time		|	C (S,I,D)			|	Go (S,I,D)
	
	t=0.01s		|	1,	290,	270		|	25,		922,	903
	N=0.1s		|	9,	2898,	2879	|	248,	19728,	19708
	N=1.00s		|	46,	28179,	28159	|	2686,	197905,	197886
	N=2.00s		|	97,	57066,	574047	|	5485,	395239,	395219

	
	