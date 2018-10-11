Chat bot code

These are the beginnings of some code for a chat bot simulation.
A user can send a command to the chat bot to leave a message for
another logged out user when they log back in.

This code is my own implementation.

Correctness
	
	Much like the linked list problem, this is basically just a pair of thread safe data
	structures which handle the chatbot and the irc channel. From an analytical perspective
	we can be sure there are no deadlocks because it is not possible for a user to 
	grab one lock and wait on a lock held by another user who is waiting on the
	first lock. Each function call is atomic from the user perspective.
	Output files are provided in "cout.txt" and "gout.txt" for C and Go respectively.
	
Comprehensibility

	This is another case where both semaphores and channels are equivelant.
	The implementation is fairly easy to read either way, as it just involves
	two mutexes on both data structures.
	
Performance

	Time		|	C messages	|	Go messages
	
	t=1.0s		|	33			|	75
	t=2.0s		|	85			|	173
	t=5.0s		|	243			|	470
	t=10.0s		|	490			|	948
	
	As usual Go's primitives can complete more operations more quickly than C's.
	Since the program is randomized and has simulated work, this does not represent much,
	but given the same logic, Go still outperforms C.
	
	