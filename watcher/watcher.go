package watcher

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

type Watcher struct {
	*websocket.Conn
	TxUpdateCh <-chan BTCTxApiEvent
	context.CancelFunc
}

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

func GetWatcher(ctx context.Context, pingInterval int) (*Watcher, error) {
	_, cancel := context.WithCancel(ctx)
	u, err := url.Parse(wsAPIAddress)
	if err != nil {
		return nil, fmt.Errorf("getWsDialer:url.Parse:[%v]", err.Error())
	}
	log.Debug().Msgf("GetWsDialer:connecting to %s", u.String())

	c, _, err := websocket.DefaultDialer.Dial(u.String(), nil)
	if err != nil {
		return nil, fmt.Errorf("getWsDialer:Dial:[%v]", err.Error())
	}
	w := &Watcher{
		Conn:       c,
		CancelFunc: cancel,
	}
	w.startPingWsAPI(ctx, pingInterval)
	err = w.SubToUnconfirmed()
	w.TxUpdateCh = w.getBtcTxUpdateChan(ctx, pingInterval)
	return w, err
}

func (w *Watcher) getBtcTxUpdateChan(ctx context.Context, pingInterval int) <-chan BTCTxApiEvent {
	out := make(chan BTCTxApiEvent)
	go func() {
		for {
			select {
			case <-ctx.Done(): // if cancel() execute
				return
			default:
			}
			_, msg, err := w.ReadMessage()
			if err != nil {
				//TODO: reconnect on error
				log.Error().Err(err).Msgf("processWsEvents:s.btcWsApi.ReadMessage:[%s]", err.Error())
				return
			}
			// fmt.Printf("msg --- %s", msg)
			if len(msg) != 0 {
				event := &BTCTxApiEvent{}
				err = json.Unmarshal(msg, event)
				if err != nil {
					log.Error().Err(err).Msgf("processWsEvents:json.Unmarshal:[%s] msg[%s]", err.Error(), msg)
				}
				out <- *event
			}
		}
	}()
	return out
}

func (w *Watcher) startPingWsAPI(ctx context.Context, pingInterval int) {
	go func() {
		for {
			select {
			case <-ctx.Done(): // if cancel() execute
				return
			case <-time.After(time.Duration(pingInterval) * time.Second):
				err := w.pingWs()
				if err != nil {
					log.Error().Err(err).Msgf("pingWsApi:pingWs:[%s]", err.Error())
				}
				break
			}
		}
	}()
}

func (w *Watcher) pingWs() error {
	return w.WriteMessage(websocket.TextMessage,
		[]byte(`{"op":"ping"}`))
}

func (w *Watcher) SubToAddress(address string) error {
	return w.WriteMessage(websocket.TextMessage,
		[]byte(fmt.Sprintf(`{"op":"addr_sub", "addr":"%s"}`, address)))
}

func (w *Watcher) UnsubFromAddress(address string) error {
	return w.WriteMessage(websocket.TextMessage,
		[]byte(fmt.Sprintf(`{"op":"addr_unsub", "addr":"%s"}`, address)))
}

func (w *Watcher) SubToUnconfirmed() error {
	return w.WriteMessage(websocket.TextMessage,
		[]byte(`{"op":"unconfirmed_sub"}`))
}
