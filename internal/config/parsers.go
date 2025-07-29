package config

import (
	"flag"
	"fmt"
	"os"
	"strconv"
	"strings"
)

const trimChars string = "/:,. "

func ParseServerFlags() {
	flag.StringVar(&ServerHost, "host", "localhost", "Server run host")
	flag.IntVar(&ServerPort, "port", 8080, "Server run port")
	flag.StringVar(&ServerEndpoint, "endpoint", "http://localhost:8080/", "Server endpoint basic url")
	flag.Parse()

	if host := os.Getenv("SERVER_HOST"); host != "" {
		ServerHost = host
	}
	if port := os.Getenv("SERVER_PORT"); port != "" {
		port, err := strconv.Atoi(port)
		if err == nil {
			ServerPort = port
		}
	}
	if endpoint := os.Getenv("SERVER_ENDPOINT"); endpoint != "" {
		ServerEndpoint = endpoint
	}

	ServerEndpoint = strings.Trim(ServerEndpoint, trimChars) + "/"
}

func ParseClientFlags() {
	flag.StringVar(&ClientEndpoint, "endpoint", "http://localhost:8080/", "Client endpoint basic url")
	flag.Parse()

	if endpoint := os.Getenv("CLIENT_ENDPOINT"); endpoint != "" {
		fmt.Println("ECHO", endpoint)
		ClientEndpoint = endpoint
	}

	ClientEndpoint = strings.Trim(ClientEndpoint, trimChars) + "/"
}
