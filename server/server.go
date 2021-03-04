package server

import (
	"context"
	"encoding/json"
	"fmt"
	"sync"

	"github.com/go-resty/resty/v2"
	"github.com/gorilla/websocket"
	"github.com/vsergeev/btckeygenie/btckey"
	"github.com/vsergeev/btckeygenie/request"
	bolt "go.etcd.io/bbolt"
)

type Server struct {
	*Config
	*request.HTTPClient
	*bolt.DB
	btcWsApi      *websocket.Conn
	pool          *Keys
	paymentsInfo  *PaymentsInfo
	btcTxUpdateCh <-chan BTCTxApiEvent
}

type Config struct {
	Port         string `env:"SERVER_PORT" envDefault:"8080"`
	XAPIKey      string `env:"X_API_KEY" envDefault:"kek"`
	AddressTTL   int    `env:"ADDRESS_TTL" envDefault:"15"` // min
	DBPath       string `env:"DB_PATH" envDefault:"payments.db"`
	KeysBucket   string `env:"KEYS_BUCKET" envDefault:"keys"`
	OrdersBucket string `env:"ORDERS_BUCKET" envDefault:"orders"`
	Debug        bool   `env:"DEBUG" envDefault:"true"`
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
	btcws, err := getWsDialer()
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
		paymentsInfo: &PaymentsInfo{
			m:     map[string]PaymentInfo{},
			Mutex: &sync.Mutex{},
		},
		btcWsApi: btcws,
	}
	btcCh := s.getBtcTxUpdateChan(context.Background())
	s.btcTxUpdateCh = btcCh
	go s.processIncoming()
	return s, nil
}
