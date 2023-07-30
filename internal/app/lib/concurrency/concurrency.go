// package concurrency содержит код реализации конкурентной обработки идентификаторов
package concurrency

import (
	"fmt"
	"strings"
	"sync"
)

// URLDeleter определяет метод SetDeleteFlag
type URLDeleter interface {
	SetDeleteFlag(IDs []string, userID string) error
}

// ProcessIDs конкурентно обрабатывает идентификаторы сокращённых ссылок
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
	_ = urlDeleter.SetDeleteFlag(urls, userID)
}
