package config

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestShortenerURL(t *testing.T) {
	testCases := []struct {
		nameTest string
		baseAddr string
		runAddr  string
	}{
		{nameTest: "#1 test", baseAddr: "http://localhost:8080/", runAddr: "localhost:8080"},
	}

	for _, tc := range testCases {
		t.Run(tc.nameTest, func(t *testing.T) {
			flags := ParseFlags()

			assert.Equal(t, tc.baseAddr, flags.FlagBaseAddr, "Флаг базового адресса сервера, не совпадает с адрессом сервера")
			assert.Equal(t, tc.runAddr, flags.FlagRunAddr, "Флаг адресса запуска сервера, не совпадает с адрессом сервера")
		})
	}
}
