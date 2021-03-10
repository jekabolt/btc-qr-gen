package server

import (
	"context"
	"encoding/json"
	"fmt"
	"sync"

	"github.com/go-resty/resty/v2"
	"github.com/vsergeev/btckeygenie/btckey"
	"github.com/vsergeev/btckeygenie/order"
	"github.com/vsergeev/btckeygenie/request"
	"github.com/vsergeev/btckeygenie/store"
	"github.com/vsergeev/btckeygenie/watcher"
)

type Server struct {
	*Config
	*request.HTTPClient
	*store.DB
	*watcher.Watcher
	pool         *Keys
	paymentsInfo *order.PaymentsInfo
}

type Config struct {
	Port           string `env:"SERVER_PORT" envDefault:"8080"`
	XAPIKey        string `env:"X_API_KEY" envDefault:"kek"`
	AddressTTL     int    `env:"ADDRESS_TTL" envDefault:"15"`      // min
	PingIntervalWS int    `env:"PING_INTERVAL_WS" envDefault:"10"` // sec
	DBPath         string `env:"DB_PATH" envDefault:"payments.db"`
	KeysBucket     string `env:"KEYS_BUCKET" envDefault:"keys"`
	OrdersBucket   string `env:"ORDERS_BUCKET" envDefault:"orders"`
	Debug          bool   `env:"DEBUG" envDefault:"true"`
}

func (c *Config) String() string {
	bs, _ := json.Marshal(c)
	return string(bs)
}

func (c *Config) InitServer() (*Server, error) {
	db, err := store.InitDB(c.DBPath, c.KeysBucket, c.OrdersBucket)
	if err != nil {
		return nil, fmt.Errorf("InitServer:GetWsDialer [%v]", err.Error())
	}
	w, err := watcher.GetWatcher(context.Background(), c.PingIntervalWS)
	if err != nil {
		return nil, fmt.Errorf("InitServer:GetWsDialer [%v]", err.Error())
	}

	r := resty.New()
	s := &Server{
		Config:     c,
		DB:         db,
		HTTPClient: request.NewHTTPClient(r),
		pool: &Keys{
			m:     map[string]btckey.BTCKeyPair{},
			Mutex: &sync.Mutex{},
		},
		paymentsInfo: order.NewPaymentsInfo(),
		Watcher:      w,
	}
	// s.SubToUnconfirmed()
	go s.processIncoming()
	return s, nil
}
