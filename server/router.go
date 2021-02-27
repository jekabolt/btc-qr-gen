package server

import (
	"log"
	"net/http"
	"time"

	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
)

func (s *Server) Serve() error {
	r := chi.NewRouter()

	r.Use(middleware.RequestID)
	r.Use(middleware.RealIP)
	r.Use(middleware.Recoverer)
	r.Use(middleware.Timeout(60 * time.Second))

	r.HandleFunc("/", s.healthCheck)
	r.Options("/*", handleOptions)

	r.Route("/v1", func(r chi.Router) {
		r.Use(middleware.Logger)
		r.Use(s.xAPICheckMiddleware)
		r.Get("/qr/{amount}/{meta}", s.getAddressQrCode)
	})

	log.Println("Listening on :" + s.Port)
	return http.ListenAndServe(":"+s.Port, r)

}
