package server

import (
	"context"
	"fmt"
	"time"

	"github.com/rs/zerolog/log"
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
	err = s.StoreBTCKeyPair(btckp)
	if err != nil {
		return nil, fmt.Errorf("getAddress:storeBTCKeyPair[%v]", err.Error())
	}
	return btckp, nil
}

func (s *Server) watchAddress(ctx context.Context, kp *btckey.BTCKeyPair) {
	go func() {
		err := s.SubToAddress(kp.AddressCompressed)
		if err != nil {
			log.Error().Err(err).Msgf("watchAddress:subToAddress:[%s]", err.Error())
		}
		for {
			select {
			case <-ctx.Done(): // if cancel() execute
				s.UnsubFromAddress(kp.AddressCompressed)
				s.pool.add(*kp)
				return
			case <-time.After(time.Duration(s.AddressTTL) * time.Minute):
				s.UnsubFromAddress(kp.AddressCompressed)
				s.pool.add(*kp)
				break
			}
		}

	}()
}

func (s *Server) processIncoming() {
	// calculate amount for outputs and update info on pool of addresses and orders. mark order as payed in db and send email
	for tx := range s.TxUpdateCh {
		for _, out := range tx.Transaction.Out {
			pi, ok := s.paymentsInfo.Get(out.Addr)
			if !ok {
				continue
			}
			pi.Success = true
			pi.Txid = tx.Transaction.Hash
			s.paymentsInfo.Add(pi)
			err := s.StorePaymentInfo(pi)
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
