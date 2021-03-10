package watcher

import (
	"context"
	"fmt"
	"testing"

	"github.com/gorilla/websocket"
)

func TestConnect(t *testing.T) {
	w, err := GetWatcher(context.Background(), 10)
	if err != nil {
		t.Fatalf("getWsDialer %s", err.Error())
	}

	ch := w.getBtcTxUpdateChan(context.Background(), 10)

	err = w.WriteMessage(websocket.TextMessage, []byte(`{"op":"unconfirmed_sub"}`))
	if err != nil {
		t.Fatalf("WriteMessage %s", err.Error())
	}

	for v := range ch {
		fmt.Println("--  ", v)
	}
}
