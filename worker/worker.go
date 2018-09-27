// Worker is module which is focused on task processing (it's like a thread in threadpool)
// Each worker registers in pool and waits for job. When it gets a job, it executes it and returns result.
package worker

type Job func()

// Worker represents the worker that executes the job
type Worker struct {
	workerPool	chan chan Job
	jobChannel	chan Job
	exit		chan bool
}

// NewWorker creates worker prepared to be started
func NewWorker(workerPool chan chan Job) Worker {
	return Worker{
		workerPool: workerPool,
		jobChannel: make(chan Job),
		exit:	   make(chan bool)}
}


// Start -- start the main loop for a worker, listening for work and a quit signal
func (w Worker) Start() {
	go func() {
		for {
			// register the current worker into the worker queue.
			w.workerPool <- w.jobChannel

			select {
			case job := <-w.jobChannel:
				job()
			case <-w.exit:
				// We have received a signal to stop
				return
			}
		}
	}()
}


// Stop signals the worker to stop listening for work requests.
func (w Worker) Stop() {
	go func() {
		w.exit <- true
	}()
}

