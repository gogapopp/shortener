package userurls

import (
	"fmt"
	"net/http"
)

// ExampleGetPingDBHandler example of working with GetURLsHandler
func ExampleGetURLsHandler() {
	c := http.Client{}
	req, err := http.NewRequest(http.MethodDelete, "localhost:8080/api/user/urls", nil)
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