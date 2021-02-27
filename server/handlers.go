package server

import (
	"fmt"
	"net/http"

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
	btcKp, err := btckey.GenerateBTCKeyPair()
	if err != nil {
		log.Error().Err(err).Msgf("getAddressQrCode:GenerateBTCKeyPair:[%s]", err.Error())
		w.Write([]byte(http.StatusText(http.StatusInternalServerError)))
		w.WriteHeader(http.StatusInternalServerError)
	}
	// TODO: save
	message := chi.URLParam(r, "meta")
	amount := chi.URLParam(r, "amount")
	png, err := qrcode.Encode(fmt.Sprintf("bitcoin:%s?amount=%s&message=%s", btcKp.AddressCompressed, amount, message), qrcode.Medium, 256)
	if err != nil {
		log.Error().Err(err).Msgf("getAddressQrCode:qrcode.Encode:[%s]", err.Error())
		w.Write([]byte(http.StatusText(http.StatusInternalServerError)))
		w.WriteHeader(http.StatusInternalServerError)
	}
	writeImage(w, png)
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
