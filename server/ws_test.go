package server

import (
	"context"
	"fmt"
	"testing"

	"github.com/gorilla/websocket"
)

func TestConnect(t *testing.T) {
	c, err := getWsDialer()
	if err != nil {
		t.Fatalf("getWsDialer %s", err.Error())
	}
	s := Server{
		btcWsApi: c,
	}
	ch := s.getBtcTxUpdateChan(context.Background())

	err = c.WriteMessage(websocket.TextMessage, []byte(`{"op":"unconfirmed_sub"}`))
	if err != nil {
		t.Fatalf("WriteMessage %s", err.Error())
	}

	for v := range ch {
		fmt.Println("--  ", v)
	}
}
