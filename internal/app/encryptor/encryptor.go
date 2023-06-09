package encryptor

import (
	"fmt"
	"strings"

	"github.com/google/uuid"
)

// хранит flags.FlagBaseAddr из main.go
var baseAddr string

// GetBaseAddr принимает значение flags.FlagBaseAddr из main.go и сохраняет в локальной переменной baseAddr
func GetBaseAddr(str string) {
	baseAddr = str
}

// ShortenerURL функция "сжимает" строку и возрващает айди
func ShortenerURL(url string) string {
	// получаем рандомный айди
	id := uuid.New()
	address := baseAddr
	// проверяем соответсвует ли строка формату http://example.ru/
	if !strings.HasPrefix(baseAddr, "http://") {
		address = fmt.Sprintf("http://%s", address)
	}
	if !strings.HasSuffix(baseAddr, "/") {
		address = fmt.Sprintf("%s/", address)
	}
	return address + id.String()
}
