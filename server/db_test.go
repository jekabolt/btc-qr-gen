package server

import (
	"fmt"
	"testing"
	"time"

	"github.com/vsergeev/btckeygenie/btckey"
	bolt "go.etcd.io/bbolt"
)

func TestDB(t *testing.T) {

	// bss, err := json.Marshal(PaymentInfo{})

	// fmt.Printf("%s", bss)
	// return

	btckp, err := btckey.GenerateBTCKeyPair()
	if err != nil {
		t.Errorf("TestDB:GenerateBTCKeyPair [%v]", err.Error())
	}
	db, err := bolt.Open("../payments.db", 0600, nil)
	if err != nil {
		t.Errorf("TestDB:bolt.Open [%v]", err.Error())
	}
	s := Server{
		DB: db,
		Config: &Config{
			KeysBucket: "keys",
		},
	}
	err = s.updateDB([]byte(s.KeysBucket),
		[]byte(btckp.AddressCompressed),
		&KeyPair{
			BTCKeyPair:     btckp,
			InitiationTime: time.Now().Unix(),
			Payed:          false,
		})
	if err != nil {
		t.Errorf("TestDB:updateDB [%v]", err.Error())
	}

	bs, err := s.queryDB([]byte(s.KeysBucket), []byte(btckp.AddressCompressed))
	if err != nil {
		t.Errorf("TestDB:queryDB [%v]", err.Error())
	}
	fmt.Printf("1 =%s\n\n\n", bs)

	err = s.deleteKey([]byte(s.KeysBucket), []byte(btckp.AddressCompressed))
	if err != nil {
		t.Errorf("TestDB:deleteKey [%v]", err.Error())
	}

	bs, err = s.queryDB([]byte(s.KeysBucket), []byte(btckp.AddressCompressed))
	if err == nil {
		t.Errorf("TestDB:queryDB should be not found")
	}
}
