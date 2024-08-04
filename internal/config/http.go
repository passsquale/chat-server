package config

import (
	"errors"
	"net"
	"os"
)

const (
	httpHostEnvName = "HTTP_HOST"
	httpPortEnvName = "HTTP_PORT"
)

type httpConfig struct {
	host string
	port string
}

func NewHTTPConfig() (HTTPConfig, error) {
	host := os.Getenv(httpHostEnvName)
	if len(host) == 0 {
		return nil, errors.New("http host not found in environments")
	}

	port := os.Getenv(httpPortEnvName)
	if len(host) == 0 {
		return nil, errors.New("http port not found in environments")
	}

	return httpConfig{
		host: host,
		port: port,
	}, nil
}

func (cfg httpConfig) Address() string {
	return net.JoinHostPort(cfg.host, cfg.port)
}
