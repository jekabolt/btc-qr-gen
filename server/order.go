package server

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
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
