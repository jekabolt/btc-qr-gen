package server

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"sync"

	"github.com/vsergeev/btckeygenie/btckey"

	"encoding/base64"
)

type PaymentInfo struct {
	Country     string `json:"country"`
	City        string `json:"city"`
	Address     string `json:"address"`
	ZIP         string `json:"zip"`
	FullName    string `json:"fullName"`
	Phone       string `json:"phone"`
	Email       string `json:"email"`
	TotalAmount string `json:"totalAmount"`
	Success     bool   `json:"success"`
}

func decodePaymentInfo(paymentInfoB64 string) (*PaymentInfo, error) {
	piDec, err := base64.StdEncoding.DecodeString(paymentInfoB64)
	if err != nil {
		return nil, fmt.Errorf("decodePaymentInfo:base64.StdEncoding.DecodeString")
	}
	pi := &PaymentInfo{}
	err = json.Unmarshal(piDec, pi)
	return pi, err
}

type Keys struct {
	m map[string]*btckey.BTCKeyPair
	*sync.Mutex
}

func (k *Keys) add(kp *btckey.BTCKeyPair) {
	k.Lock()
	defer k.Unlock()
	k.m[kp.AddressCompressed] = kp
}

func (k *Keys) delete(address string) {
	k.Lock()
	defer k.Unlock()
	delete(k.m, address)
}

func (k *Keys) get(address string) *btckey.BTCKeyPair {
	return k.m[address]
}

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
