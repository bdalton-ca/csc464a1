#include "stdlib.h"
#include "stdio.h"
#include "string.h"
#include "time.h"
#include <pthread.h>
#include <semaphore.h>


double get_ms(){ return 1000.0f * (double)clock() / (double)CLOCKS_PER_SEC ; }

#define N_CUSTOMER 800
#define QUEUE_MAX 100
int queue_size = 0;
int newest_customer = 0;
int oldest_customer = 0;

sem_t queue[QUEUE_MAX];
sem_t qmutex;

pthread_mutex_t pfmutex;

#define printf_p(...)
//#define printf_p(...) do\
{\
	pthread_mutex_lock(&pfmutex);\
	printf(__VA_ARGS__);\
	pthread_mutex_unlock(&pfmutex);\
}while(0)\

void * customer(void*)
{	
	//printf_p("customer enters shop\n");
	for(;;)
	{		
		sem_wait(&qmutex);		
		if(queue_size < QUEUE_MAX)
		{
			printf_p("customer enters queue (waiting=%d)\n", queue_size+1, newest_customer);
		
			int id = newest_customer;
			
			newest_customer = (newest_customer+1)%QUEUE_MAX;
			queue_size++;			
			
			sem_post(&qmutex);
			sem_wait(&queue[id]);
	
			
			queue_size--;			
			
			sem_post(&qmutex);
			break;
		}
		else sem_post(&qmutex);	
	}
	printf_p("customer leaves shop\n");
}

void * barber(void*)
{
	printf_p("barber enters shop\n");
	for(int serviced = 0; serviced < N_CUSTOMER;)
	{
		sem_wait(&qmutex);
		if(queue_size > 0)
		{
			printf_p("barber services customer (waiting=%d)\n", queue_size-1);
			
			int id = oldest_customer;
			oldest_customer = (oldest_customer+1)%QUEUE_MAX;			
			
			serviced++;
			sem_post(&queue[id]);
			// customer unlocks queue on wake
		}
		else sem_post(&qmutex);
	}
	printf_p("barber exits shop\n");
}


int main(int argc, char * argv[])
{
	double t0 = get_ms();
	sem_init(&qmutex,0,1);
	pthread_mutex_init(&pfmutex,NULL);
	pthread_t threads[N_CUSTOMER];
	for(int i = 0; i < QUEUE_MAX; i++)
		sem_init(&queue[i],0,0);
	
	for(int i = 0; i < N_CUSTOMER; i++)
		pthread_create(&threads[i], NULL, customer, NULL);
	
	barber(NULL);
	
    for (int i = 0; i < N_CUSTOMER; i++)
       pthread_join(threads[i], NULL);
	
	double t1 = get_ms();
	
	printf("duration: %f",t1-t0);
	
	return 0;
}