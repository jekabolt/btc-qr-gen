package server

import (
	"context"
	"encoding/json"
	"fmt"
	"net/url"
	"time"

	"github.com/gorilla/websocket"
	"github.com/rs/zerolog/log"
	"github.com/vsergeev/btckeygenie/btckey"
)

// https://www.blockchain.com/api/api_websocket

const (
	wsAPIAddress = "wss://ws.blockchain.info/inv"
)

type BTCTxApiEvent struct {
	Operation   string `json:"op"`
	Transaction Tx     `json:"x"`
}

type Inputs struct {
	Sequence int64  `json:"sequence"`
	Script   string `json:"script"`
}
type Out struct {
	Spent   bool   `json:"spent"`
	TxIndex int    `json:"tx_index"`
	Type    int    `json:"type"`
	Addr    string `json:"addr"`
	Value   int    `json:"value"`
	N       int    `json:"n"`
	Script  string `json:"script"`
}
type Tx struct {
	Inputs []Inputs `json:"inputs"`
	Out    []Out    `json:"out"`
	Time   int      `json:"time"`
	Hash   string   `json:"hash"`
}

func getWsDialer() (*websocket.Conn, error) {
	u, err := url.Parse(wsAPIAddress)
	if err != nil {
		return nil, err
	}
	log.Debug().Msgf("GetWsDialer:connecting to %s", u.String())

	c, _, err := websocket.DefaultDialer.Dial(u.String(), nil)
	if err != nil {
		return nil, fmt.Errorf("getWsDialer:Dial:[%v]", err.Error())
	}
	return c, err
}

func (s *Server) getBtcTxUpdateChan(ctx context.Context) <-chan BTCTxApiEvent {
	s.startPingWsAPI(ctx)
	out := make(chan BTCTxApiEvent)
	go func() {
		for {
			select {
			case <-ctx.Done(): // if cancel() execute
				return
			default:
			}
			_, msg, err := s.btcWsApi.ReadMessage()
			if err != nil {
				log.Error().Err(err).Msgf("processWsEvents:s.btcWsApi.ReadMessage:[%s]", err.Error())
			}
			event := &BTCTxApiEvent{}
			err = json.Unmarshal(msg, event)
			if err != nil {
				log.Error().Err(err).Msgf("processWsEvents:json.Unmarshal:[%s] msg[%s]", err.Error(), msg)
				//TODO: reconnect on error
			}
			out <- *event
		}
	}()
	return out
}

func (s *Server) processIncoming() {
	// calculate amount for outputs and update info on pool of addresses and orders. mark order as payed in db and send email

	for tx := range s.btcTxUpdateCh {
		for _, out := range tx.Transaction.Out {
			pi, ok := s.paymentsInfo.get(out.Addr)
			if !ok {
				continue
			}
			pi.Success = true
			pi.Txid = tx.Transaction.Hash
			s.paymentsInfo.add(pi)
			err := s.storePaymentInfo(pi)
			if err != nil {
				log.Error().Err(err).Msgf("processIncoming:storePaymentInfo:[%s]", err.Error())
				continue
			}

			// TODO:
			//send ws event on tx in mempool
			//send email to payment info recipient
		}

	}
}

func (s *Server) startPingWsAPI(ctx context.Context) {
	go func() {
		for {
			select {
			case <-ctx.Done(): // if cancel() execute
				return
			case <-time.After(time.Duration(s.PingWsAPI) * time.Second):
				err := s.pingWs()
				if err != nil {
					log.Error().Err(err).Msgf("pingWsApi:pingWs:[%s]", err.Error())
				}
				break
			}
		}
	}()
}

func (s *Server) pingWs() error {
	return s.btcWsApi.WriteMessage(websocket.TextMessage,
		[]byte(`{"op":"ping"}`))
}

func (s *Server) subToAddress(address string) error {
	return s.btcWsApi.WriteMessage(websocket.TextMessage,
		[]byte(fmt.Sprintf(`{"op":"addr_sub", "addr":"%s"}`, address)))
}

func (s *Server) unsubFromAddress(address string) error {
	return s.btcWsApi.WriteMessage(websocket.TextMessage,
		[]byte(fmt.Sprintf(`{"op":"addr_unsub", "addr":"%s"}`, address)))
}

func (s *Server) subToUnconfirmed() error {
	return s.btcWsApi.WriteMessage(websocket.TextMessage,
		[]byte(`{"op":"unconfirmed_sub"}`))
}

func (s *Server) watchAddress(kp *btckey.BTCKeyPair) (context.Context, context.CancelFunc) {
	ctx, cancel := context.WithCancel(context.Background())
	go func() {
		err := s.subToAddress(kp.AddressCompressed)
		if err != nil {
			log.Error().Err(err).Msgf("watchAddress:subToAddress:[%s]", err.Error())
		}
		for {
			select {
			case <-ctx.Done(): // if cancel() execute
				s.unsubFromAddress(kp.AddressCompressed)
				s.pool.add(*kp)
				return
			case <-time.After(time.Duration(s.AddressTTL) * time.Minute):
				s.unsubFromAddress(kp.AddressCompressed)
				s.pool.add(*kp)
				break
			}
		}

	}()
	return ctx, cancel
}
