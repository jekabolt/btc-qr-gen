package server

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"sync"
)

type PaymentInfo struct {
	OrderId     string `json:"orderId"`
	BTCAddress  string `json:"btcAddress"`
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

type PaymentsInfo struct {
	m map[string]PaymentInfo // key BTCAddress
	*sync.Mutex
}

func (pi *PaymentsInfo) len() int {
	return len(pi.m)
}

func (pi *PaymentsInfo) add(info PaymentInfo) {
	pi.Lock()
	defer pi.Unlock()
	pi.m[info.BTCAddress] = info
}

func (pi *PaymentsInfo) delete(address string) {
	pi.Lock()
	defer pi.Unlock()
	delete(pi.m, address)
}

func (pi *PaymentsInfo) get(address string) PaymentInfo {
	return pi.m[address]
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
