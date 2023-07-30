package ping

import (
	"fmt"
	"net/http"
)

// ExampleGetPingDBHandler пример работы с GetPingDBHandler
func ExampleGetPingDBHandler() {
	c := http.Client{}
	req, err := http.NewRequest(http.MethodDelete, "localhost:8080/ping", nil)
	if err != nil {
		panic(err)
	}
	resp, err := c.Do(req)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()
	fmt.Println(resp.Status)
}
