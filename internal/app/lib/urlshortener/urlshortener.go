// package urlshortener contains an implementation of string shortening
package urlshortener

import (
	"crypto/rand"
	"fmt"
	"log"
	"strings"
)

// ShortenerURL function "compresses" the string and returns a short link
func ShortenerURL(baseAddr string) string {
	const size = 6
	b := make([]byte, size)
	_, err := rand.Read(b)
	if err != nil {
		log.Fatal(err)
	}
	var letters = []rune("ABCDEFGHIJKLMNOPQRSTUVWXYZ" +
		"abcdefghijklmnopqrstuvwxyz" +
		"0123456789")

	address := baseAddr
	// check if the format string matches http://example.ru/
	if !strings.HasPrefix(baseAddr, "http://") {
		address = fmt.Sprintf("http://%s", address)
	}
	if !strings.HasSuffix(baseAddr, "/") {
		address = fmt.Sprintf("%s/", address)
	}
	var result strings.Builder
	result.Write([]byte(address))
	for _, b := range b {
		result.WriteRune(letters[int(b)%len(letters)])
	}
	return result.String()
}
