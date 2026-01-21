Exercise 1 - Theory questions
-----------------------------

### Concepts

What is the difference between *concurrency* and *parallelism*?
> *concurrency is when one thread switches between different tasks to emulate parallellism. Parallellism is when several threads each execute tasks at the same time.*

What is the difference between a *race condition* and a *data race*? 
> *a data race happens when several threads access a spot in memory simultaneously. A race condition is when a value is different from what the program expects due to an unexpected order of execution.* 
 
*Very* roughly - what does a *scheduler* do, and how does it do it?
> *a scheduler recieves requests for threads that should execute, and decides where and when to execute these.* 


### Engineering

Why would we use multiple threads? What kinds of problems do threads solve?
> *multiple threads make the program safer. They can also serve as a more direct translation of a parallell design specification into programming language.*

Some languages support "fibers" (sometimes called "green threads") or "coroutines"? What are they, and why would we rather use them over threads?
> *green threads are managed by a program-specific scheduler. These have less overhead than the OS scheduler, and can function in environments without an OS scheduler.*

Does creating concurrent programs make the programmer's life easier? Harder? Maybe both?
> *for small systems, concurrency is generally simpler. For larger parallell systems parallellism may simlify understanding and implementation.*

What do you think is best - *shared variables* or *message passing*?
> *it depends. For large data shared variables are better. For small data message passing avoids a lot of race conditions.*


