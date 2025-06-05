package main

import (
	"context"
	"flag"
	"fmt"
	"sync"
	"time"
)

type WorkerPool struct {
	jobs    chan string
	wg      sync.WaitGroup
	mu      sync.Mutex
	workers map[int]context.CancelFunc
	nextID  int
}

func NewWorkerPool(size int) *WorkerPool {
	return &WorkerPool{
		jobs:    make(chan string, size),
		workers: make(map[int]context.CancelFunc),
		nextID:  1,
	}
}

func (wp *WorkerPool) worker(ctx context.Context, id int) {
	defer wp.wg.Done()

	fmt.Printf("Worker %d started\n", id)
	defer fmt.Printf("Worker %d stopped\n", id)

	for {
		select {
		case <-ctx.Done():
			return
		case task, ok := <-wp.jobs:
			if !ok {
				return
			}
			fmt.Printf("Worker %v started job %v\n", id, task)
			time.Sleep(time.Second)
			fmt.Printf("Worker %v finished job %v\n", id, task)
		}
	}
}

func (wp *WorkerPool) AddWorker() int {
	wp.mu.Lock()
	defer wp.mu.Unlock()

	id := wp.nextID
	wp.nextID++

	ctx, cancel := context.WithCancel(context.Background())
	wp.workers[id] = cancel

	wp.wg.Add(1)
	go wp.worker(ctx, id)
	return id
}

func (wp *WorkerPool) RemoveWorker(id int) bool {
	wp.mu.Lock()
	defer wp.mu.Unlock()

	cancel, ok := wp.workers[id]
	if !ok {
		return false
	}

	cancel()
	delete(wp.workers, id)
	return true
}

func (wp *WorkerPool) Close() {
	wp.mu.Lock()
	for _, cancel := range wp.workers {
		cancel()
	}
	wp.mu.Unlock()

	close(wp.jobs)
	wp.wg.Wait()
}

func (wp *WorkerPool) WorkerCount() int {
	wp.mu.Lock()
	defer wp.mu.Unlock()

	return len(wp.workers)
}

func main() {
	numWorkers := flag.Int("workers", 10, "Number of workers in the pool")
	numTasks := flag.Int("tasks", 50, "Number of tasks to process")
	flag.Parse()

	pool := NewWorkerPool(*numWorkers)

	for i := 1; i <= *numWorkers; i++ {
		pool.AddWorker()
	}
	fmt.Printf("Starting with %d workers\n", pool.WorkerCount())

	time.AfterFunc(time.Second * 3, func() {
		id := pool.AddWorker()
		fmt.Printf("Added worker %d, total workers: %d\n", id, pool.WorkerCount())
	})

	time.AfterFunc(time.Second * 5, func() {
		if pool.WorkerCount() > 0 {
			if pool.RemoveWorker(1) {
				fmt.Printf("Removed worker 1, total workers: %d\n", pool.WorkerCount())
			}
		}
	})

	for i := 1; i <= *numTasks; i++ {
		pool.jobs <- fmt.Sprintf("job-%d", i)
		time.Sleep(time.Millisecond * 500)
	}

	pool.Close()
	fmt.Println("Finish")
}
