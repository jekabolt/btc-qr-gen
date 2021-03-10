package server

import (
	"context"
	"fmt"
	"net/http"

	"github.com/go-chi/chi"
	"github.com/rs/zerolog/log"
	"github.com/skip2/go-qrcode"
	"github.com/vsergeev/btckeygenie/order"
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
	var err error

	amount, err = convertAmount(amount)
	if err != nil {
		log.Error().Err(err).Msgf("getAddressQrCode:convertAmount:[%s]", err.Error())
		writeBadRequest(w)
		return
	}
	log.Debug().Msgf("amount [%v]", amount)

	pi, err := order.DecodePaymentInfo(meta)
	if err != nil {
		log.Error().Err(err).Msgf("getAddressQrCode:decodePaymentInfo:[%s]", err.Error())
		writeBadRequest(w)
		return
	}
	s.paymentsInfo.Add(pi)

	btckp, err := s.getAddress()
	if err != nil {
		log.Error().Err(err).Msgf("getAddressQrCode:addPaymentInfo:[%s]", err.Error())
		writeInternalServerError(w)
		return
	}
	s.watchAddress(context.Background(), btckp)

	qrData := fmt.Sprintf("bitcoin:%s?amount=%s&message=%s", btckp.AddressCompressed, amount, meta)
	png, err := qrcode.Encode(qrData, qrcode.Medium, 256)
	if err != nil {
		log.Error().Err(err).Msgf("getAddressQrCode:qrcode.Encode:[%s]", err.Error())
		writeInternalServerError(w)
		return
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
