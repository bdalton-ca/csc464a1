#include "stdlib.h"
#include "stdio.h"
#include "string.h"
#include "time.h"
#include <pthread.h>
#include <semaphore.h>

double get_ms(){ return 1000.0f * (double)clock() / (double)CLOCKS_PER_SEC ; }

#define N_CUSTOMER 400
#define PARTY_MAX 5
int queue_size = 0;
int party_size = 0;
int is_party = 0;

sem_t mutex;
sem_t block;

pthread_mutex_t pfmutex;
//#define printf_p(...)
#define printf_p(...) do\
{\
	pthread_mutex_lock(&pfmutex);\
	printf(__VA_ARGS__);\
	pthread_mutex_unlock(&pfmutex);\
}while(0)\

void * customer(void*)
{	
	sem_wait(&mutex);
	if(is_party)
	{
		queue_size++;
		sem_post(&mutex);
		sem_wait(&block);
		queue_size--;
	}

	party_size++;
	printf_p("customer joins table (party_size=%d)\n",party_size);
	
	if(party_size == PARTY_MAX) 
	{
		is_party = 1;
		printf_p("There is a full party of customers\n",party_size);
	}
	
	if(!is_party && queue_size>0)
		sem_post(&block);
	else
		sem_post(&mutex);
	
	printf_p("customer eats sushi\n");
	
	sem_wait(&mutex);
	party_size--;
	if(party_size == 0) 
	{
		is_party = 0;
		printf_p("party leaves table\n",party_size);
	}
	else if(!is_party)
		printf_p("customer leaves table (party_size=%d)\n",party_size);
	
	if(!is_party && queue_size>0)
		sem_post(&block);
	else
		sem_post(&mutex);	
}

int main(int argc, char * argv[])
{
	double t0 = get_ms();
	
	pthread_mutex_init(&pfmutex,NULL);
	pthread_t threads[N_CUSTOMER];
	
	sem_init(&mutex,0,1);
	sem_init(&block,0,0);
	
	double t1 = get_ms();
	
	for(int i = 0; i < N_CUSTOMER; i++)
	{
		pthread_create(&threads[i], NULL, customer, NULL);
	}
	
    for (int i = 0; i < N_CUSTOMER; i++)
       pthread_join(threads[i], NULL);
	
	double t2 = get_ms();
	
	printf("init time: %f\n", t1-t0);
	printf("proc time: %f\n", t2-t1);
	printf("total: %f\n", t2-t0);
		
	return 0;
}