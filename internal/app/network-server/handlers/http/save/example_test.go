package save

import (
	"fmt"
	"io"
	"net/http"
	"strings"
)

// ExampleGetPingDBHandler example of working with PostSaveHandler
func ExamplePostSaveHandler() {
	request := "https://practicum.yandex.ru"
	resp, err := http.Post("localhost:8080", "text/plain", strings.NewReader(request))
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		panic(err)
	}
	fmt.Print(string(body))
}
