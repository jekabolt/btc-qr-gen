package server

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/go-chi/chi"
	"github.com/rs/zerolog/log"
	"github.com/skip2/go-qrcode"
	"github.com/vsergeev/btckeygenie/btckey"
)

func (s *Server) healthCheck(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("OK"))
}

func setCORSHeaders(w http.ResponseWriter) {
	w.Header().Set("Access-Control-Allow-Methods", "GET,POST,OPTIONS")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Headers", "*")
}

func handleOptions(w http.ResponseWriter, r *http.Request) {
	setCORSHeaders(w)
}

func (s *Server) getAddressQrCode(w http.ResponseWriter, r *http.Request) {
	meta := chi.URLParam(r, "meta")
	amount := chi.URLParam(r, "amount")
	pi, err := decodePaymentInfo(meta)
	if err != nil {
		log.Error().Err(err).Msgf("getAddressQrCode:decodePaymentInfo:[%s]", err.Error())
		writeBadRequest(w)
	}
	btckp, err := btckey.GenerateBTCKeyPair()
	if err != nil {
		log.Error().Err(err).Msgf("getAddressQrCode:GenerateBTCKeyPair:[%s]", err.Error())
		writeInternalServerError(w)
	}
	err = s.updateDB([]byte(s.KeysBucket),
		[]byte(btckp.AddressCompressed),
		&KeyPair{
			BTCKeyPair:     btckp,
			InitiationTime: time.Now().Unix(),
			Payed:          false,
		})
	if err != nil {
		log.Error().Err(err).Msgf("getAddressQrCode:addKeyPair:[%s]", err.Error())

	}
	err = s.updateDB([]byte(s.OrdersBucket),
		[]byte(btckp.AddressCompressed), pi)
	if err != nil {
		log.Error().Err(err).Msgf("getAddressQrCode:addPaymentInfo:[%s]", err.Error())
		writeInternalServerError(w)
	}

	ctx, _ := context.WithCancel(context.Background())
	// TODO: cancel
	s.watchAddress(ctx, btckp.AddressCompressed)

	qrData := fmt.Sprintf("bitcoin:%s?amount=%s&message=%s", btckp.AddressCompressed, amount, meta)
	png, err := qrcode.Encode(qrData, qrcode.Medium, 256)
	if err != nil {
		log.Error().Err(err).Msgf("getAddressQrCode:qrcode.Encode:[%s]", err.Error())
		writeInternalServerError(w)
	}
	err = writeImage(w, png)
	if err != nil {
		log.Error().Err(err).Msgf("getAddressQrCode:writeImage:[%s]", err.Error())
	}
}

func (s *Server) xAPICheckMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		headerKey := r.Header.Get("x-api-key")
		queryKey := r.URL.Query().Get("x-api-key")
		if s.XAPIKey != headerKey && s.XAPIKey != queryKey {
			w.WriteHeader(http.StatusUnauthorized)
			w.Write([]byte("bad x-api-key"))
			return
		}
		next.ServeHTTP(w, r)
	})
}
