// Package utils provides utility functions for the goprofit application.
package utils

import "sync"

// WorkerPool manages a pool of workers for parallel processing with controlled concurrency.
type WorkerPool struct {
	sem chan struct{}
	wg  sync.WaitGroup
}

// NewWorkerPool creates a new pool with n concurrent workers.
func NewWorkerPool(n int) *WorkerPool {
	return &WorkerPool{
		sem: make(chan struct{}, n),
	}
}

// Submit adds a task to the pool. It blocks if all workers are busy.
func (wp *WorkerPool) Submit(task func()) {
	wp.wg.Add(1)
	wp.sem <- struct{}{} // Acquire slot (blocks if pool is full)
	go func() {
		defer wp.wg.Done()
		defer func() { <-wp.sem }() // Release slot
		task()
	}()
}

// Wait blocks until all submitted tasks complete.
func (wp *WorkerPool) Wait() {
	wp.wg.Wait()
}

// Close releases resources. Should be called when the pool is no longer needed.
func (wp *WorkerPool) Close() {
	close(wp.sem)
}
