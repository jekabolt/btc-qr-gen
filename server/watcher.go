package server

import (
	"context"
	"time"
)

// https://www.blockchain.com/api/api_websocket

func subToAddress(address string) error {
	return nil
}

func unsubFromAddress(address string) error {
	return nil
}

func (s *Server) watchAddress(ctx context.Context, address string) {
	go func() {
		err := subToAddress(address)
		if err != nil {
			// return fmt.Errorf("watchAddress:subToAddress[%v]", err.Error())
		}
		for {
			select {
			case <-ctx.Done(): // if cancel() execute
				unsubFromAddress(address)
			case <-time.After(time.Duration(s.AddressTTL) * time.Minute):
				unsubFromAddress(address)
			}
		}
	}()
}
