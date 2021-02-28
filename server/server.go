package server

import (
	"encoding/json"
	"fmt"

	"gitlab.com/miapago/report-server/request"
	bolt "go.etcd.io/bbolt"
)

type Server struct {
	*Config
	*request.HTTPClient
	*bolt.DB
}

type Config struct {
	Port       string `env:"SERVER_PORT" envDefault:"8080"`
	XAPIKey    string `env:"X_API_KEY" envDefault:"kek"`
	AddressTTL string `env:"ADDRESS_TTL" envDefault:"kek"`
	DBPath     string `env:"DB_PATH" envDefault:"keys.db"`
	KeysBucket string `env:"KEYS_BUCKET" envDefault:"keys"`
	Debug      bool   `env:"DEBUG" envDefault:"true"`
}

func (c *Config) String() string {
	bs, _ := json.Marshal(c)
	return string(bs)
}

func (c *Config) InitServer() (*Server, error) {
	db, err := bolt.Open(c.DBPath, 0600, nil)
	if err != nil {
		return nil, fmt.Errorf("InitServer:bolt.Open [%v]", err.Error())
	}
	s := &Server{
		Config: c,
		DB:     db,
	}
	s.createBucket(c.KeysBucket)
	return s, nil
}
