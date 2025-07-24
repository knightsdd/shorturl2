package config

type mainConfig struct {
	runAddress    runAddress
	prefixAddress prefixAddress
}

type otherConfig struct {
	endpoint prefixAddress
}

var serverConfig mainConfig
var clientConfig otherConfig

func GetServerRunAddress() string {
	return serverConfig.runAddress.String()
}

func GetServerPrefixAddress() string {
	return serverConfig.prefixAddress.String()
}

func GetClientEndpoint() string {
	return clientConfig.endpoint.String()
}
