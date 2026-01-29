package worker

import (
	"context"
	"fmt"
	"log"
	"sync"
	"time"
)

// Job represents a background job
type Job struct {
	ID        string
	Type      string
	Payload   interface{}
	CreatedAt time.Time
	Status    string
	Error     error
}

// JobHandler defines job execution interface
type JobHandler interface {
	Handle(ctx context.Context, job *Job) error
}

// Worker executes background jobs
type Worker struct {
	id       int
	jobQueue chan *Job
	handlers map[string]JobHandler
	quit     chan bool
	wg       *sync.WaitGroup
}

// WorkerPool manages multiple workers
type WorkerPool struct {
	workers   []*Worker
	jobQueue  chan *Job
	handlers  map[string]JobHandler
	quit      chan bool
	wg        sync.WaitGroup
	isRunning bool
	mu        sync.Mutex
}

// NewWorkerPool creates a new worker pool
func NewWorkerPool(workerCount int) *WorkerPool {
	return &WorkerPool{
		workers:  make([]*Worker, workerCount),
		jobQueue: make(chan *Job, 100),
		handlers: make(map[string]JobHandler),
		quit:     make(chan bool),
	}
}

// RegisterHandler registers a job handler
func (wp *WorkerPool) RegisterHandler(jobType string, handler JobHandler) {
	wp.handlers[jobType] = handler
}

// Start starts the worker pool
func (wp *WorkerPool) Start() {
	wp.mu.Lock()
	defer wp.mu.Unlock()

	if wp.isRunning {
		return
	}

	for i := 0; i < len(wp.workers); i++ {
		worker := &Worker{
			id:       i + 1,
			jobQueue: wp.jobQueue,
			handlers: wp.handlers,
			quit:     make(chan bool),
			wg:       &wp.wg,
		}
		wp.workers[i] = worker
		wp.wg.Add(1)
		go worker.start()
	}

	wp.isRunning = true
	log.Printf("Worker pool started with %d workers", len(wp.workers))
}

// Stop stops the worker pool
func (wp *WorkerPool) Stop() {
	wp.mu.Lock()
	defer wp.mu.Unlock()

	if !wp.isRunning {
		return
	}

	close(wp.quit)
	for _, worker := range wp.workers {
		worker.stop()
	}
	wp.wg.Wait()
	close(wp.jobQueue)
	wp.isRunning = false
	log.Println("Worker pool stopped")
}

// Submit submits a job to the worker pool
func (wp *WorkerPool) Submit(job *Job) error {
	wp.mu.Lock()
	defer wp.mu.Unlock()

	if !wp.isRunning {
		return fmt.Errorf("worker pool is not running")
	}

	select {
	case wp.jobQueue <- job:
		return nil
	default:
		return fmt.Errorf("job queue is full")
	}
}

// start starts the worker
func (w *Worker) start() {
	defer w.wg.Done()

	for {
		select {
		case job := <-w.jobQueue:
			w.executeJob(job)
		case <-w.quit:
			return
		}
	}
}

// stop stops the worker
func (w *Worker) stop() {
	close(w.quit)
}

// executeJob executes a job
func (w *Worker) executeJob(job *Job) {
	log.Printf("Worker %d: Processing job %s (type: %s)", w.id, job.ID, job.Type)

	handler, exists := w.handlers[job.Type]
	if !exists {
		log.Printf("Worker %d: No handler for job type %s", w.id, job.Type)
		job.Status = "failed"
		job.Error = fmt.Errorf("no handler for job type: %s", job.Type)
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Minute)
	defer cancel()

	if err := handler.Handle(ctx, job); err != nil {
		log.Printf("Worker %d: Job %s failed: %v", w.id, job.ID, err)
		job.Status = "failed"
		job.Error = err
	} else {
		log.Printf("Worker %d: Job %s completed successfully", w.id, job.ID)
		job.Status = "completed"
	}
}

// Example job handlers

// EmailJobHandler handles email sending jobs
type EmailJobHandler struct{}

func (h *EmailJobHandler) Handle(ctx context.Context, job *Job) error {
	payload, ok := job.Payload.(map[string]interface{})
	if !ok {
		return fmt.Errorf("invalid payload format")
	}

	to := payload["to"].(string)
	subject := payload["subject"].(string)
	body := payload["body"].(string)

	log.Printf("Sending email to %s: %s", to, subject)

	// Simulate email sending
	time.Sleep(2 * time.Second)

	// In production, use actual email service (SendGrid, AWS SES, etc.)
	log.Printf("Email sent successfully: %s", body)

	return nil
}

// DataProcessingJobHandler handles data processing jobs
type DataProcessingJobHandler struct{}

func (h *DataProcessingJobHandler) Handle(ctx context.Context, job *Job) error {
	payload, ok := job.Payload.(map[string]interface{})
	if !ok {
		return fmt.Errorf("invalid payload format")
	}

	log.Printf("Processing data: %v", payload)

	// Simulate data processing
	time.Sleep(3 * time.Second)

	log.Println("Data processing completed")

	return nil
}

// UserCleanupJobHandler handles user cleanup jobs
type UserCleanupJobHandler struct{}

func (h *UserCleanupJobHandler) Handle(ctx context.Context, job *Job) error {
	log.Println("Running user cleanup job...")

	// Simulate cleanup
	time.Sleep(1 * time.Second)

	log.Println("User cleanup completed")

	return nil
}
