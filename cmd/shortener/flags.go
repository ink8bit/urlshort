package main

import (
	"flag"
)

var addr string

func parseFlags() {
	var baseURL string

	flag.StringVar(&addr, "a", ":8080", "address and port to run server")
	flag.StringVar(&baseURL, "b", "", "short url prefix")
	flag.Parse()
}
