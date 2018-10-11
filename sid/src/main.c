#include "stdlib.h"
#include "stdio.h"
#include "string.h"
#include "time.h"
#include <pthread.h>
#include <semaphore.h>
#include <windows.h>

void sleep(float ms) { Sleep((unsigned)ms); }

#define N_NODES 20
#define N_THREADS 3

typedef struct NODE_t
{
	int data;
	struct NODE_t * next;
} NODE;

NODE * list;
NODE * free_nodes;
int list_size=0;

int search_count=0;
int insert_count=0;

sem_t search_mutex;
sem_t insert_mutex;
sem_t delete_mutex;
sem_t no_searchers;
sem_t no_inserters;

int n_search=0;
int n_insert=0;
int n_delete=0;

sem_t pfmutex;
sem_t rnmutex;
//#define printf_p(...)
#define printf_p(...) do\
{\
	sem_wait(&pfmutex);\
	printf(__VA_ARGS__);\
	sem_post(&pfmutex);\
}while(0)\

unsigned rand32()
{
	unsigned x;
	#define RAND_SH(X) ((((unsigned)rand())&0xFF)<<X)
	sem_wait(&rnmutex);\
	x = RAND_SH(0) | RAND_SH(8) | RAND_SH(16) | RAND_SH(24);
	sem_post(&rnmutex);\
	return x;
}


void * searcher(void*)
{
	for(;;)
	{
		//if((rand32()&0xFFF)>0) continue;
		sem_wait(&search_mutex); 
		if(search_count==0) sem_wait(&no_searchers);
		search_count++;
		sem_post(&search_mutex);
		
		int size = list_size;
		
		if(size > 0)
		{
			int r = rand32()%size;
			NODE * node = list;
			for( int i = 0; node->next && i < r; i++ ) node=node->next;
			printf_p("searcher found %d, (size=%d)\n",node->data, size);
		}
		
		sem_wait(&search_mutex); 
		n_search++;
		search_count--;
		if(search_count==0) sem_post(&no_searchers);
		sem_post(&search_mutex);
		
	}
}

void * inserter(void*)
{	
	for(;;)
	{
		sem_wait(&no_inserters);
		if(N_NODES-list_size > 0)
		{
			int r = list_size? rand32()%list_size : 0;
			NODE * node = list;
			for( int i = 0; i < r; i++ ) node=node->next;
			
			NODE * tmp = free_nodes;
			free_nodes = free_nodes->next;
			
			if(node)
			{
				tmp->next = node->next;
				node->next = tmp;
			}
			else
			{
				tmp->next = NULL;
				list = tmp;
			}
			
			list_size++;
			
			
			printf_p("inserter added %d, (size=%d)\n",tmp->data,list_size);
			
			sem_post(&insert_mutex);
			n_insert++;
			sem_wait(&insert_mutex);
		}		
		sem_post(&no_inserters);
	}	
}

void * deleter(void*)
{
	for(;;)
	{		
		sem_wait(&no_searchers);
		sem_wait(&no_inserters);
		if(list_size > 0)
		{
			int r = list_size-1? rand32()%(list_size-1) : 0;
			NODE * node = list;
			for( int i = 0; i < r; i++ ) node=node->next;
			
			NODE * tmp = node->next? node->next : node;
			node->next = node->next? node->next->next : NULL;
			tmp->next = free_nodes;
			free_nodes = tmp;
			

			printf_p("deleter removed %d\n",tmp->data);
			
			list_size--;
			sem_post(&delete_mutex);
			n_delete++;
			sem_wait(&delete_mutex);
		}
		sem_post(&no_searchers);
		sem_post(&no_inserters);
	}
	
}

int main(int argc, char * argv[])
{	
	sem_init(&search_mutex,0,1);
	sem_init(&insert_mutex,0,1);
	sem_init(&delete_mutex,0,1);
	sem_init(&no_searchers,0,1);
	sem_init(&no_inserters,0,1);
	sem_init(&pfmutex,0,1);
	sem_init(&rnmutex,0,1);
	
	
	NODE nodes[N_NODES];
	for(int i = 0; i < N_NODES; i++)
	{
		nodes[i].data = i;
		nodes[i].next = &nodes[i+1];
	}
	nodes[N_NODES-1].next = NULL;
	free_nodes = nodes;
	
	
	pthread_t threads[N_THREADS*3];
	for(int i = 0; i < N_THREADS*3; i+=3)
	{
		
		pthread_create(&threads[i+0], NULL, searcher, NULL);
		pthread_create(&threads[i+1], NULL, inserter, NULL);
		pthread_create(&threads[i+2], NULL, deleter, NULL);
	}
	
	sleep(1000.0f);
	
	//printf("sid: %d,%d,%d",n_search,n_insert,n_delete);
	
	for(;;);
	
	return 0;
}