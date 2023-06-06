package concurrency

import (
	"sync"

	"github.com/gogapopp/shortener/internal/app/auth"
)

// Allocate заполняет буфер
func Allocate(jobs chan string, IDs []string) {
	for _, id := range IDs {
		jobs <- id
	}
	close(jobs)
}

// Worker читает буфер и вызывает метод SetDeleteFlag
func Worker(wg *sync.WaitGroup, jobs chan string, userID string, baseAddr string) {
	for job := range jobs {
		auth.GlobalStore.SetDeleteFlag(userID, job, baseAddr)
	}
	wg.Done()
}

// CreateWorkerPool создаёт пул горутин и ожидает их завершения
func CreateWorkerPool(noOfWorkers int, jobs chan string, userID string, baseAddr string) {
	var wg sync.WaitGroup
	for i := 0; i < noOfWorkers; i++ {
		wg.Add(1)
		go Worker(&wg, jobs, userID, baseAddr)
	}
	wg.Wait()
}
