package config

import (
	"flag"
	"strings"
)

func ParseServerFlags() {
	flag.StringVar(&ServerHost, "host", "localhost", "Server run host")
	flag.IntVar(&ServerPort, "port", 8080, "Server run port")
	flag.StringVar(&ServerEndpoint, "endpoint", "http://localhost:8080", "Server endpoint basic url")
	flag.Parse()

	ServerEndpoint = strings.Trim(ServerEndpoint, "/:,. ") + "/"
}

func ParseClientFlags() {
	flag.StringVar(&ClientEndpoint, "endpoint", "http://localhost:8080", "Client endpoint basic url")
	flag.Parse()
	ClientEndpoint = strings.Trim(ClientEndpoint, "/:,. ") + "/"
}
