package config

import (
	"flag"
)

type Flags struct {
	FlagRunAddr     string
	FlagBaseAddr    string
	FlagStoragePath string
}

// ParseFlags() парсит флаги при запуске программы
func ParseFlags() *Flags {
	f := &Flags{}
	flag.StringVar(&f.FlagRunAddr, "a", "localhost:8080", "address and port to run server")
	flag.StringVar(&f.FlagBaseAddr, "b", "http://localhost:8080/", "base url")
	flag.StringVar(&f.FlagStoragePath, "f", "", "file storage path")
	flag.Parse()
	return f
}
