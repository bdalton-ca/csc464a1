#include "stdlib.h"
#include "stdio.h"
#include "string.h"
#include "time.h"
#include <pthread.h>
#include <semaphore.h>

double get_ms(){ return 1000.0f * (double)clock() / (double)CLOCKS_PER_SEC ; }

#define CAR_CAP 20
#define N_PASSENGERS 100

sem_t mutex_queue;
sem_t mutex_depart;
sem_t mutex_unboard;
sem_t mutex_ready;
int passenger_count=0;
int unboarded_count=0;


pthread_mutex_t pfmutex;
//#define printf_p(...)
#define printf_p(...) do\
{\
	pthread_mutex_lock(&pfmutex);\
	printf(__VA_ARGS__);\
	pthread_mutex_unlock(&pfmutex);\
}while(0)\

void * car(void*)
{
	printf_p("Ready to load\n");
	sem_post(&mutex_queue);
	for(int i = 0; i < N_PASSENGERS;) 
	{
		sem_wait(&mutex_queue);
		int count = passenger_count;		
		if (count == CAR_CAP) 
		{			
			printf_p("Car is running\n");
			sem_wait(&mutex_depart);
			printf_p("Car is unloading\n");
			for(int i = 0; i < CAR_CAP; i++) sem_post(&mutex_unboard);
			for(int i = 0; i < CAR_CAP; i++) sem_wait(&mutex_ready);
			printf_p("Ready to load\n");
			count = 0;
			i+= CAR_CAP;
		}
		passenger_count = count;
		sem_post(&mutex_queue);
	}
}

void * passenger(void*)
{
	for(;;)
	{
		sem_wait(&mutex_queue);
		int count = passenger_count;
		if (count < CAR_CAP)
		{
			count += 1;
			printf_p("%d passengers have boarded,\n", count);
			passenger_count = count;
			if(count==CAR_CAP) sem_post(&mutex_depart);
			sem_post(&mutex_queue);
			
			sem_wait(&mutex_unboard);
			printf_p("passenger unboarded\n");
			sem_post(&mutex_ready);
			break;
			/*
			count = unboarded_count;
			count -= 1;
			printf_p("%d passenger have unboarded\n", CAR_CAP-count);
			unboarded_count = count;
			sem_post(&mutex_unboard);
			*/
		} 
		else 
		{
			passenger_count = count;
			sem_post(&mutex_queue);
		}
	}
}


int main(int argc, char * argv[])
{
	double t0 = get_ms();
	
	sem_init(&mutex_queue,0,0);
	sem_init(&mutex_unboard,0,0);
	sem_init(&mutex_depart,0,0);
	sem_init(&mutex_ready,0,0);
	pthread_mutex_init(&pfmutex,NULL);
	
	//pthread_mutex_lock(&mutex_unboard);
	//pthread_mutex_lock(&mutex_queue);
	
	pthread_t threads[N_PASSENGERS];
	for(int i = 0; i < N_PASSENGERS; i++)	
		pthread_create(&threads[i], NULL, passenger, NULL);
	
	double t1 = get_ms();
	
	car(NULL);
	
	double t2 = get_ms();
	
	printf("init time: %f\n", t1-t0);
	printf("proc time: %f\n", t2-t1);
	
	return 0;
}