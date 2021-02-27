package server

import (
	"encoding/json"

	"gitlab.com/miapago/report-server/request"
)

type Server struct {
	*Config
	*request.HTTPClient
}

type Config struct {
	Port    string `env:"SERVER_PORT" envDefault:"8080"`
	XAPIKey string `env:"X_API_KEY" envDefault:"kek"`
	Debug   bool   `env:"DEBUG" envDefault:"true"`
}

func (c *Config) String() string {
	bs, _ := json.Marshal(c)
	return string(bs)
}

func (c *Config) InitServer() (*Server, error) {
	s := &Server{
		Config: c,
	}
	return s, nil
}
