package server

import (
	"fmt"

	"github.com/vsergeev/btckeygenie/btckey"
)

func (s *Server) getAddress() (*btckey.BTCKeyPair, error) {
	if s.pool.len() != 0 {
		for _, kp := range s.pool.m {
			return &kp, nil
		}
	}
	btckp, err := btckey.GenerateBTCKeyPair()
	if err != nil {
		return nil, fmt.Errorf("getAddress:GenerateBTCKeyPair[%v]", err.Error())
	}
	err = s.storeBTCKeyPair(btckp)
	if err != nil {
		return nil, fmt.Errorf("getAddress:storeBTCKeyPair[%v]", err.Error())
	}
	return btckp, nil
}
