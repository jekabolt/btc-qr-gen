package store

import (
	"fmt"
	"testing"
	"time"

	"github.com/vsergeev/btckeygenie/btckey"
)

func TestDB(t *testing.T) {

	// bss, err := json.Marshal(PaymentInfo{})

	// fmt.Printf("%s", bss)
	// return

	btckp, err := btckey.GenerateBTCKeyPair()
	if err != nil {
		t.Errorf("TestDB:GenerateBTCKeyPair [%v]", err.Error())
	}

	db, err := InitDB("../payment.db", "keys", "orders")
	if err != nil {
		t.Errorf("TestDB:InitDB [%v]", err.Error())
	}

	err = db.updateDB([]byte(db.KeysBucket),
		[]byte(btckp.AddressCompressed),
		&KeyPair{
			BTCKeyPair:     btckp,
			InitiationTime: time.Now().Unix(),
			Payed:          false,
		})
	if err != nil {
		t.Errorf("TestDB:updateDB [%v]", err.Error())
	}

	bs, err := db.queryDB([]byte(db.KeysBucket), []byte(btckp.AddressCompressed))
	if err != nil {
		t.Errorf("TestDB:queryDB [%v]", err.Error())
	}
	fmt.Printf("1 =%s\n\n\n", bs)

	err = db.deleteKey([]byte(db.KeysBucket), []byte(btckp.AddressCompressed))
	if err != nil {
		t.Errorf("TestDB:deleteKey [%v]", err.Error())
	}

	bs, err = db.queryDB([]byte(db.KeysBucket), []byte(btckp.AddressCompressed))
	if err == nil {
		t.Errorf("TestDB:queryDB should be not found")
	}
}
