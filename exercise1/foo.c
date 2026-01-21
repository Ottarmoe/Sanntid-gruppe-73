// Compile with `gcc foo.c -Wall -std=gnu99 -lpthread`, or use the makefile
// The executable will be named `foo` if you use the makefile, or `a.out` if you use gcc directly

#include <pthread.h>
#include <stdio.h>

int i = 0;
pthread_mutex_t mut;

// Note the return type: void*
void* incrementingThreadFunction(){
    // TODO: increment i 1_000_000 times
    for(int j = 0; j<1000000; j++) {
        pthread_mutex_lock(&mut);
        i++;
        pthread_mutex_unlock(&mut);
    }
    printf("inc: %d\n", i);
    return NULL;
}

void* decrementingThreadFunction(){
    // TODO: decrement i 1_000_000 times
    for(int j = 0; j<1000000; j++) {
        pthread_mutex_lock(&mut);
        i--;
        pthread_mutex_unlock(&mut);
    }
    printf("dec: %d\n", i);
    return NULL;
}


int main(){
    // TODO: 
    // start the two functions as their own threads using `pthread_create`
    // Hint: search the web! Maybe try "pthread_create example"?
    pthread_mutex_init(&mut, NULL);
    pthread_t a, b;
    pthread_create(&a, NULL, incrementingThreadFunction, NULL);
    pthread_create(&b, NULL, decrementingThreadFunction, NULL);

    pthread_join(a, NULL);
    pthread_join(b, NULL);
    // TODO:
    // wait for the two threads to be done before printing the final result
    // Hint: Use `pthread_join`    
    
    printf("The magic number is: %d\n", i);
    return 0;
}