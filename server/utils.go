package server

import (
	"bytes"
	"net/http"
	"strconv"
)

func writeImage(w http.ResponseWriter, img []byte) error {
	bytes.NewBuffer(img)
	buffer := bytes.NewBuffer(img)
	w.Header().Set("Content-Type", "image/png")
	w.Header().Set("Content-Length", strconv.Itoa(len(buffer.Bytes())))
	_, err := w.Write(buffer.Bytes())
	return err
}

func writeInternalServerError(w http.ResponseWriter) {
	w.Write([]byte(http.StatusText(http.StatusInternalServerError)))
	w.WriteHeader(http.StatusInternalServerError)
}

func writeBadRequest(w http.ResponseWriter) {
	w.Write([]byte(http.StatusText(http.StatusBadRequest)))
	w.WriteHeader(http.StatusBadRequest)
}
