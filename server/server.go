package server

import (
	"encoding/json"
	"fmt"

	"github.com/go-resty/resty/v2"
	"github.com/vsergeev/btckeygenie/request"
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
	AddressTTL int    `env:"ADDRESS_TTL" envDefault:"10"` // min
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
	r := resty.New()
	s := &Server{
		Config:     c,
		DB:         db,
		HTTPClient: request.NewHTTPClient(r),
	}
	return s, nil
}
