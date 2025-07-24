package config

import (
	"errors"
	"flag"
	"fmt"
	"strconv"
	"strings"
)

type runAddress struct {
	host string
	port int
}

func (na *runAddress) String() string {
	return fmt.Sprintf("%s:%d", na.host, na.port)
}

func (na *runAddress) Set(flagValue string) error {
	values := strings.Split(flagValue, ":")
	if len(values) != 2 {
		return errors.New("need address ia a format host:port")
	}
	na.host = values[0]
	port, err := strconv.Atoi(values[1])
	if err != nil {
		return fmt.Errorf("port must be int value, input '%s': %w", values[1], err)
	}
	na.port = port
	return nil
}

type prefixAddress struct {
	protocol string
	host     string
	port     int
}

func (ba *prefixAddress) String() string {
	return fmt.Sprintf("%s://%s:%d/", ba.protocol, ba.host, ba.port)
}

func (ba *prefixAddress) Set(flagValue string) error {
	values := strings.Split(flagValue, ":")
	if len(values) != 3 {
		return errors.New("need address ia a format protocol://host:port")
	}
	port, err := strconv.Atoi(strings.Trim(values[2], "/"))
	if err != nil {
		return fmt.Errorf("port must be int value, input '%s': %w", values[1], err)
	}
	ba.protocol = strings.Trim(values[0], "/")
	ba.host = strings.Trim(values[1], "/")
	ba.port = port
	return nil
}

func ParseServerFlags() {
	runAddr := &runAddress{
		host: "localhost",
		port: 8080,
	}
	prefixAddr := new(prefixAddress)

	flag.Var(runAddr, "a", "Server run address host:port")
	flag.Var(prefixAddr, "b", "Base address for requests protokol://host:port")

	flag.Parse()

	if prefixAddr.protocol == "" {
		prefixAddr.protocol = "http"
		prefixAddr.host = runAddr.host
		prefixAddr.port = runAddr.port
	}

	serverConfig.runAddress = *runAddr
	serverConfig.prefixAddress = *prefixAddr
}

func ParseClientFlags() {
	endpointAddr := &prefixAddress{
		protocol: "http",
		host:     "localhost",
		port:     8080,
	}

	flag.Var(endpointAddr, "a", "Address to send requests protocol:://host:port")

	flag.Parse()

	clientConfig.endpoint = *endpointAddr
}
