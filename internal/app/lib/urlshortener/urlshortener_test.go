package urlshortener

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestShortenerURL(t *testing.T) {
	baseAddr := "https://www.example.com"
	shortURLs := make(map[string]string)

	// we fill the map with 10 links, in the loop we check the uniqueness of each
	for i := 0; i < 10; i++ {
		shortURL := ShortenerURL(baseAddr)
		assert.False(t, strings.HasPrefix(shortURL, baseAddr+"/"))
		// checking if there is such a shortURL
		_, ok := shortURLs[shortURL]
		assert.False(t, ok)
		shortURLs[shortURL] = shortURL
	}
}
