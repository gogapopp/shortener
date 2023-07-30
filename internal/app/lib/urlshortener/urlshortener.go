package urlshortener

import (
	"fmt"
	"math/rand"
	"strings"
	"time"
)

// ShortenerURL функция "сжимает" строку и возрващает короткую ссылку
func ShortenerURL(baseAddr string) string {
	const size = 6
	rnd := rand.New(rand.NewSource(time.Now().UnixNano()))

	chars := []rune("ABCDEFGHIJKLMNOPQRSTUVWXYZ" +
		"abcdefghijklmnopqrstuvwxyz" +
		"0123456789")

	b := make([]rune, size)
	for i := range b {
		b[i] = chars[rnd.Intn(len(chars))]
	}

	address := baseAddr
	// проверяем соответсвует ли строка формату http://example.ru/
	if !strings.HasPrefix(baseAddr, "http://") {
		address = fmt.Sprintf("http://%s", address)
	}
	if !strings.HasSuffix(baseAddr, "/") {
		address = fmt.Sprintf("%s/", address)
	}

	return fmt.Sprintf("%s%s", address, string(b))
}
