#include "stdlib.h"
#include "stdio.h"
#include "string.h"
#include "time.h"
#include <pthread.h>
#include <semaphore.h>


double get_ms(){ return 1000.0f * (double)clock() / (double)CLOCKS_PER_SEC ; }

sem_t mutex;
pthread_mutex_t pfmutex;

#define printf_p(...)
//#define printf_p(...) do\
{\
	pthread_mutex_lock(&pfmutex);\
	printf(__VA_ARGS__);\
	pthread_mutex_unlock(&pfmutex);\
}while(0)\

sem_t hydrogen_sp;
sem_t oxygen_sp;
pthread_barrier_t barrier;

int n_oxygen_waiting = 0;
int n_hydrogen_waiting = 0;
	
void * oxygen(void*data)
{
	//ATOM * atom = (ATOM*)data;
	sem_wait(&mutex);
	n_oxygen_waiting++;
	if(n_hydrogen_waiting>=2)
	{
		pthread_barrier_init(&barrier,NULL,3);
		sem_post(&hydrogen_sp);
		sem_post(&hydrogen_sp);
		sem_post(&oxygen_sp);
		n_hydrogen_waiting-=2;
		n_oxygen_waiting--;
	}
	else
		sem_post(&mutex);
	
	sem_wait(&oxygen_sp);
	printf_p("bonding oxygen\n");	
	
	pthread_barrier_wait(&barrier);
	pthread_barrier_destroy(&barrier);
	printf_p("FORMED H20\n");
	sem_post(&mutex);
}

void * hydrogen(void*data)
{
	//ATOM * atom = (ATOM*)data;
	sem_wait(&mutex);
	n_hydrogen_waiting++;
	if(n_hydrogen_waiting>=2 && n_oxygen_waiting>=1)
	{
		
		pthread_barrier_init(&barrier,NULL,3);
		sem_post(&hydrogen_sp);
		sem_post(&hydrogen_sp);
		sem_post(&oxygen_sp);
		n_hydrogen_waiting-=2;
		n_oxygen_waiting--;
	}
	else
		sem_post(&mutex);	
	
	sem_wait(&hydrogen_sp);
	printf_p("bonding hydrogen\n");
	
	pthread_barrier_wait(&barrier);
}

int main(int argc, char * argv[])
{
	//#define N_ATOMS 9
	//#define N_ATOMS 99
	//#define N_ATOMS 399
	#define N_ATOMS 801
	int atom_types[N_ATOMS];
	
	double t0 = get_ms();
	
	
	srand(clock());
	sem_init(&hydrogen_sp,0,0);
	sem_init(&oxygen_sp,0,0);
	sem_init(&mutex,0,1);
	pthread_mutex_init(&pfmutex,NULL);
	
	{	// init atoms
		
		// make sure we have enough of each
		// so that they can all be made into molecules
		for(int i = 0; i < N_ATOMS; i+=3)
		{
			atom_types[i+0] = 0;
			atom_types[i+1] = 0;
			atom_types[i+2] = 1;
		}
		
		// shuffle atoms
		for(int i =  0 ; i < N_ATOMS  ; i++)
		for(int j = i+1; j < N_ATOMS-1; j++)
		{
			int dst = j + rand()%(N_ATOMS-j);
			int tmp = atom_types[i];
			atom_types[i] = atom_types[dst];
			atom_types[dst] = tmp; 
		}
	}
	
	pthread_t threads[N_ATOMS];
	for(int i = 0; i < N_ATOMS; i++)
		pthread_create(&threads[i], NULL, atom_types[i]? oxygen:hydrogen, NULL);
	
    for (int i = 0; i < N_ATOMS; i++)
       pthread_join(threads[i], NULL);
   
	double t1 = get_ms();
	
	printf("duration: %f", t1-t0);
   
	return 0;
}