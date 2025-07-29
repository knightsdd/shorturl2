package config

import (
	"fmt"
)

// server config
var ServerHost string
var ServerPort int
var ServerEndpoint string

// client config
var ClientEndpoint string

// server get config funcs
func GetServerRunAddress() string {
	return fmt.Sprintf("%s:%d", ServerHost, ServerPort)
}
