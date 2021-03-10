package order

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"sync"
)

type OrderInfo struct {
	OrderId  string   `json:"orderId"`
	Country  string   `json:"country"`
	City     string   `json:"city"`
	Address  string   `json:"address"`
	ZIP      string   `json:"zip"`
	FullName string   `json:"fullName"`
	Phone    string   `json:"phone"`
	Email    string   `json:"email"`
	ItemsId  []string `json:"itemsId"`
}

type PaymentInfo struct {
	*OrderInfo
	BTCAddress  string `json:"btcAddress"`
	Txid        string `json:"txid"`
	TotalAmount string `json:"totalAmount"`
	Success     bool   `json:"success"`
}

func NewPaymentsInfo() *PaymentsInfo {
	return &PaymentsInfo{
		m:     map[string]PaymentInfo{},
		Mutex: &sync.Mutex{},
	}
}

type PaymentsInfo struct {
	m map[string]PaymentInfo // key BTCAddress
	*sync.Mutex
}

func (pi *PaymentsInfo) len() int {
	return len(pi.m)
}

func (pi *PaymentsInfo) Add(info *PaymentInfo) {
	pi.Lock()
	defer pi.Unlock()
	pi.m[info.BTCAddress] = *info
}

func (pi *PaymentsInfo) Delete(address string) {
	pi.Lock()
	defer pi.Unlock()
	delete(pi.m, address)
}

func (pi *PaymentsInfo) Get(address string) (*PaymentInfo, bool) {
	pinfo, ok := pi.m[address]
	return &pinfo, ok
}

func DecodePaymentInfo(paymentInfoB64 string) (*PaymentInfo, error) {
	piDec, err := base64.StdEncoding.DecodeString(paymentInfoB64)
	if err != nil {
		return nil, fmt.Errorf("decodePaymentInfo:base64.StdEncoding.DecodeString")
	}
	pi := &PaymentInfo{}
	err = json.Unmarshal(piDec, pi)
	return pi, err
}
