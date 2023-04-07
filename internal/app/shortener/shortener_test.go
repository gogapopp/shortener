package shortener

import (
	"strings"
	"testing"
)

var addres string

func TestShortenerURL(t *testing.T) {
	testCases := []struct {
		nameTest    string
		url         string
		expectedUrl string
	}{
		{nameTest: "#1 test", url: "localhost:8080", expectedUrl: "http://localhost:8080/"},
		{nameTest: "#2 test", url: "localhost:8080", expectedUrl: "localhost:8080"},
	}

	for _, tc := range testCases {
		t.Run(tc.nameTest, func(t *testing.T) {
			resultURL := ShortenerURL(tc.url)
			if !strings.HasPrefix(resultURL, "http://") {
				addres = "http://" + "localhost:8080"
			}
			if !strings.HasSuffix(resultURL, "/") {
				addres = "localhost:8080" + "/"
			}
		})
	}
}
