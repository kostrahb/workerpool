# Worker Pool
Simple completely generic job queue for golang. It handles only job distribution and nothing more. Based on Malwarebytes blog and their [Handling 1M requests per minute post](http://marcio.io/2015/07/handling-1-million-requests-per-minute-with-golang/).

# Example usage

As the worker pool is intended to be as simple and generic as possible, some features (e.g. waiting) is left on user.
The pool requires usage of closures which shields it completely from business logic so the usage is following:

```golang
package main

import(
	"os"
	"fmt"
	"sync"
	"strconv"
	"github.com/kostrahb/workerpool"
)

func main() {
	MaxWorker, _ := strconv.Atoi(os.Getenv("MAX_WORKERS"))
	if MaxWorker == 0 {
		MaxWorker = 4
	}

	d := workerpool.NewPool(MaxWorker)
	d.Start()

	// Worker pool does not care about waiting (for architectural reasons) so if you want to wait you have to make this mechanism yourself.
	// It isn't that hard is it? :)
	var wg sync.WaitGroup

	for i := 0; i < 100; i++ {
		// Printing i directly is not a good idea...
		// (we're accessing it from multiple goroutines concurently and it is not protected against race conditions)
		a := i
		wg.Add(1)
		work := func() {
			fmt.Println(a)
			wg.Done()
		}
		d.AddWork(work)
	}
	wg.Wait()
}
```
