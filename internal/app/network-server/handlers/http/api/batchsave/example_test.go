package batchsave

import (
	"encoding/json"
	"net/http"
	"strings"
)

// ExampleGetPingDBHandler example of working with PostBatchJSONhHandler
func ExamplePostBatchJSONhHandler() {
	request := `
	[
		{
			"correlation_id": "1",
			"original_url": "https://practicum.yandex.ru"
		},
		{
			"correlation_id": "2",
			"original_url": "https://google.com"
		}
	]`
	expect := []struct {
		CorrelationID string `json:"correlation_id"`
		ShortURL      string `json:"short_url"`
	}{}

	resp, err := http.Post("localhost:8080/api/shorten/batch", "application/json", strings.NewReader(request))
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()
	err = json.NewDecoder(resp.Body).Decode(&expect)
	if err != nil {
		panic(err)
	}
}
