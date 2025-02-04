package config

import "github.com/joho/godotenv"

func Load(path string) error {
	err := godotenv.Load(path)
	if err != nil {
		return err
	}

	return nil
}

type SwaggerConfig interface {
	Address() string
}

type HTTPConfig interface {
	Address() string
}

type GRPCConfig interface {
	Address() string
}

type PGConfig interface {
	DSN() string
}

type AuthConfig interface {
	Address() string
}
