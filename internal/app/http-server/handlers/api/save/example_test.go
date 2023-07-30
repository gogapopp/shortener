package save

import (
	"encoding/json"
	"net/http"
	"strings"
)

// ExampleGetPingDBHandler пример работы с PostSaveJSONHandler
func ExamplePostSaveJSONHandler() {
	request := `{"url":"https://practicum.yandex.ru"}`
	expect := struct {
		ShortURL string `json:"result"`
	}{}

	resp, err := http.Post("localhost:8080/api/shorten", "application/json", strings.NewReader(request))
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()
	err = json.NewDecoder(resp.Body).Decode(&expect)
	if err != nil {
		panic(err)
	}
}
