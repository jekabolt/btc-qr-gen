package server

import (
	"bytes"
	"fmt"
	"net/http"
	"strconv"

	"github.com/btcsuite/btcutil"
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

func convertAmount(amount string) (string, error) {
	a, err := strconv.Atoi(amount)
	if err != nil {
		return "", fmt.Errorf("convertAmount:strconv.Atoi[%v]", err.Error())
	}
	amount = btcutil.Amount(a).Format(btcutil.AmountBTC)
	if len(amount) > 3 {
		return amount[:len(amount)-3], nil
	}
	return "", fmt.Errorf("convertAmount:err bad amount  [%s]", amount)
}
