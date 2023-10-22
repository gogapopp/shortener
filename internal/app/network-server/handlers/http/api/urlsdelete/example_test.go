package urlsdelete

import (
	"fmt"
	"net/http"
	"strings"
)

// ExampleGetPingDBHandler example of working with DeleteHandler
func ExampleDeleteHandler() {
	request := `[
		"MWmHmO",
		"/yRxA7V",
		"/NOMtJ6",
		"/88c078e7-452d-477b-8b55-70633284c97e"
	]`
	c := http.Client{}
	req, err := http.NewRequest(http.MethodDelete, "localhost:8080/api/user/urls", strings.NewReader(request))
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
