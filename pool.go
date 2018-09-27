// Worker pool is module which is focused on work distribution to workers (it simulates a threadpool)
// Each worker pool is created with given number of workers and exposes its joq queue via function AddWork()
package workerpool


import(
	"github.com/kostrahb/workerpool/worker"
)

// Pool is main structure in this package, it holds all information needed for work distribution and worker signaling
type Pool struct {
	// Channel that we can send work requests on.
	jobQueue chan worker.Job
	// Signaling of exit
	exit chan bool
	// A pool of workers channels that are registered within the pool
	workerPool chan chan worker.Job
	workers []worker.Worker
}


// NewPool creates dispachter with given number of workers
func NewPool(workers int) *Pool {
	jobQueue := make(chan worker.Job)
	exit := make(chan bool)
	workerPool := make(chan chan worker.Job, workers)
	d := &Pool{workerPool: workerPool, exit: exit, jobQueue:jobQueue}

	for i := 0; i < workers; i++ {
		w := worker.NewWorker(d.workerPool)
		d.workers = append(d.workers, w)
	}

	return d
}


// Start all workers and listen on queue for work
func (d *Pool) Start() {
	for _, w := range(d.workers) {
		w.Start()
	}
	go d.dispatch()
}

func (d *Pool) AddWork(job worker.Job) {
	d.jobQueue <- job
}

func (d *Pool) dispatch() {
	for {
		select {
		case job := <-d.jobQueue:
			// a job request has been received
			go func(job worker.Job) {
				// try to obtain a worker job channel that is available.
				// this will block until a worker is idle
				jobChannel := <-d.workerPool

				// dispatch the job to the worker job channel
				jobChannel <- job
			}(job)
		case <-d.exit:
			return
		}
	}
}


// Stop Pool and all workers
func (d *Pool) Stop() {
	go func() {
		d.exit <- true
		for _, w := range(d.workers) {
			w.Stop()
		}
	}()
}
