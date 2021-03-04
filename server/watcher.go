package server

import (
	"context"
	"encoding/json"
	"fmt"
	"net/url"
	"time"

	"github.com/gorilla/websocket"
	"github.com/rs/zerolog/log"
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

	//init
	// err = c.WriteJSON(`{"op":"ping"}`)
	// if err != nil {
	// 	return nil, fmt.Errorf("getWsDialer:WriteJSON:[%v]", err.Error())
	// }
	return c, err
}

func (s *Server) getBtcTxUpdateChan(ctx context.Context) <-chan BTCTxApiEvent {
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
				log.Error().Err(err).Msgf("processWsEvents:json.Unmarshal:[%s]", err.Error())
			}
			out <- *event
		}
	}()
	return out
}

func (s *Server) subToAddress(address string) error {
	return s.btcWsApi.WriteMessage(websocket.TextMessage,
		[]byte(fmt.Sprintf(`{"op":"addr_sub", "addr":"%s"}`, address)))
}

func (s *Server) unsubFromAddress(address string) error {
	return s.btcWsApi.WriteMessage(websocket.TextMessage,
		[]byte(fmt.Sprintf(`{"op":"addr_unsub", "addr":"%s"}`, address)))
}

func (s *Server) watchAddress(ctx context.Context, address string) {
	go func() {
		err := s.subToAddress(address)
		if err != nil {
			// return fmt.Errorf("watchAddress:subToAddress[%v]", err.Error())
		}
		for {
			select {
			case <-ctx.Done(): // if cancel() execute
				s.unsubFromAddress(address)
			case <-time.After(time.Duration(s.AddressTTL) * time.Minute):
				s.unsubFromAddress(address)
			}
		}
	}()
}
