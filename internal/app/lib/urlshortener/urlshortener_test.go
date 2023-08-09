package urlshortener

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestShortenerURL(t *testing.T) {
	baseAddr := "https://www.example.com"
	shortURLs := make(map[string]string)

	for i := 0; i < 10; i++ {
		shortURL := ShortenerURL(baseAddr)
		assert.False(t, strings.HasPrefix(shortURL, baseAddr+"/"))
		// проверяем есть ли такой shortURL
		_, ok := shortURLs[shortURL]
		assert.False(t, ok)
		t.Log(shortURL)
		shortURLs[shortURL] = shortURL
	}
}
