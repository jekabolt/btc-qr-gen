package server

import (
	"bytes"
	"net/http"
	"strconv"

	"github.com/rs/zerolog/log"
)

func writeImage(w http.ResponseWriter, img []byte) {

	bytes.NewBuffer(img)

	buffer := bytes.NewBuffer(img)
	w.Header().Set("Content-Type", "image/png")
	w.Header().Set("Content-Length", strconv.Itoa(len(buffer.Bytes())))
	if _, err := w.Write(buffer.Bytes()); err != nil {
		log.Error().Err(err).Msgf("writeImage:w.Write[%s]", err.Error())
		w.Write([]byte(http.StatusText(http.StatusInternalServerError)))
		w.WriteHeader(http.StatusInternalServerError)
	}
}
