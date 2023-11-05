package config

import (
	"flag"
)

var (
	Addr    string
	BaseURL string
)

func ParseFlags() {
	flag.StringVar(&Addr, "a", ":8080", "address and port to run server")
	flag.StringVar(&BaseURL, "b", "", "short url prefix")
	flag.Parse()
}
