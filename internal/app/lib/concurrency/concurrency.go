// package concurrency contains code for implementing competitive id processing
package concurrency

import (
	"fmt"
	"log"
	"strings"
	"sync"
)

// URLDeleter defines the SetDeleteFlag method
type URLDeleter interface {
	SetDeleteFlag(IDs []string, userID string) error
}

// ProcessIDs competitively handles abbreviated link IDs
func ProcessIDs(IDs []string, reqURL string, urlDeleter URLDeleter, userID string) {
	idsCh := make(chan string, len(IDs)+2)
	urlsCh := make(chan string, len(IDs)+2)
	var wg sync.WaitGroup
	for _, id := range IDs {
		idsCh <- id
	}
	close(idsCh)
	for i := 0; i < 4; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for id := range idsCh {
				url := id
				if !strings.HasPrefix(id, "/") {
					url = fmt.Sprintf("%s/%s", reqURL, id)
				} else if !strings.HasPrefix(id, "http") {
					url = fmt.Sprintf("%s%s", reqURL, id)
				}
				urlsCh <- url
			}
		}()
	}
	go func() {
		wg.Wait()
		close(urlsCh)
	}()
	var urls []string
	for url := range urlsCh {
		urls = append(urls, url)
	}
	err := urlDeleter.SetDeleteFlag(urls, userID)
	if err != nil {
		log.Printf("failed set delete flag: %s", err)
	}
}
