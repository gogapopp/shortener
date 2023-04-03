package config

import "flag"

type Flags struct {
	FlagRunAddr  string
	FlagBaseAddr string
}

func ParseFlags() *Flags {
	f := &Flags{}
	flag.StringVar(&f.FlagRunAddr, "a", "localhost:8080", "address and port to run server")
	flag.StringVar(&f.FlagBaseAddr, "b", "http://localhost:8080/", "base url")
	flag.Parse()
	return f
}
