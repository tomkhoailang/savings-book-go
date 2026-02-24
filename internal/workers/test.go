package workers


import (
"fmt"
"math/rand"
"sync"
"time"
)

// Job represents a unit of work for the workers.
type Job struct {
	ID int
}

// Worker function that processes jobs from the job channel.
func worker(id int, jobs <-chan Job, results chan<- int) {
	for j := range jobs {
		fmt.Printf("Worker %d started job %d\n", id, j.ID)
		time.Sleep(time.Duration(rand.Intn(1000)) * time.Millisecond) // Simulate work
		fmt.Printf("Worker %d finished job %d\n", id, j.ID)
		results <- j.ID
	}
}

func main() {
	rand.Seed(time.Now().UnixNano())

	// Create channels for communication.
	jobs := make(chan Job, 100) // Buffered channel to hold jobs
	results := make(chan int, 100)

	// Start workers
	numWorkers := 3
	var wg sync.WaitGroup
	for w := 1; w <= numWorkers; w++ {
		wg.Add(1)
		go func(w int) {
			defer wg.Done()
			worker(w, jobs, results)
		}(w)
	}

	// Send jobs to workers
	numJobs := 10
	for j := 1; j <= numJobs; j++ {
		jobs <- Job{ID: j}
	}
	close(jobs) // Signal that no more jobs will be sent

	// Collect results from workers
	go func() {
		wg.Wait()
		close(results) // Signal that all results have been received
	}()

	// Print results as they arrive
	for r := range results {
		fmt.Printf("Received result: %d\n", r)
	}

	fmt.Println("All jobs completed!")
}