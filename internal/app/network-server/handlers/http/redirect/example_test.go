package redirect

import (
	"fmt"
	"net/http"
)

// ExampleGetPingDBHandler example of working with GetURLGetterHandler
func ExampleGetURLGetterHandler() {
	c := http.Client{}
	req, err := http.NewRequest(http.MethodDelete, "localhost:8080/{id}", nil)
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
