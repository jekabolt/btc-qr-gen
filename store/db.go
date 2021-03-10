package store

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/vsergeev/btckeygenie/btckey"
	"github.com/vsergeev/btckeygenie/order"
	bolt "go.etcd.io/bbolt"
)

type DB struct {
	*bolt.DB
	DBPath       string
	KeysBucket   string
	OrdersBucket string
}

func InitDB(DBPath, keysBucket, ordersBucket string) (*DB, error) {
	db, err := bolt.Open(DBPath, 0600, nil)
	if err != nil {
		return nil, fmt.Errorf("InitServer:bolt.Open [%v]", err.Error())
	}
	return &DB{
		DB:           db,
		KeysBucket:   keysBucket,
		OrdersBucket: ordersBucket,
	}, nil
}

type KeyPair struct {
	*btckey.BTCKeyPair
	InitiationTime int64 `json:"initiationTime,omitempty"`
	Payed          bool  `json:"payed,omitempty"`
}

func (db *DB) updateDB(bucketName, key []byte, value interface{}) error {
	bs, err := json.Marshal(value)
	if err != nil {
		return fmt.Errorf("updateDB:json.Marsha:%v", err.Error())
	}
	return db.Update(func(tx *bolt.Tx) error {
		bkt, err := tx.CreateBucketIfNotExists(bucketName)
		if err != nil {
			return err
		}
		err = bkt.Put(key, bs)
		if err != nil {
			return err
		}
		return nil
	})
}

func (db *DB) queryDB(bucketName, key []byte) ([]byte, error) {
	v := []byte{}
	err := db.View(func(tx *bolt.Tx) error {
		bkt := tx.Bucket(bucketName)
		if bkt == nil {
			return fmt.Errorf("bucket %s not found", bucketName)
		}
		v = bkt.Get(key)
		return nil
	})
	if len(v) == 0 {
		return nil, fmt.Errorf("queryDB:notfound")
	}
	return v, err
}

func (db *DB) iterateDB(bucketName []byte) error {
	return db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket(bucketName)
		c := b.Cursor()
		for k, v := c.First(); k != nil; k, v = c.Next() {
			fmt.Printf("k =[%s], v=[%s]\n", k, v)
		}
		return nil
	})
}

func (db *DB) deleteKey(bucketName, keyName []byte) error {
	return db.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket(bucketName)
		err := b.Delete(keyName)
		return err
	})
}

func (db *DB) StoreBTCKeyPair(kp *btckey.BTCKeyPair) error {
	err := db.updateDB([]byte(db.KeysBucket),
		[]byte(kp.AddressCompressed),
		&KeyPair{
			BTCKeyPair:     kp,
			InitiationTime: time.Now().Unix(),
			Payed:          false,
		})
	if err != nil {
		return fmt.Errorf("storeOrderInfo:s.updateDB:storeBTCKeyPair[%v]", err.Error())
	}
	return nil
}

func (db *DB) StorePaymentInfo(pi *order.PaymentInfo) error {
	err := db.updateDB([]byte(db.OrdersBucket),
		[]byte(pi.BTCAddress),
		pi)
	if err != nil {
		return fmt.Errorf("storeOrderInfo:s.updateDB:storeBTCKeyPair[%v]", err.Error())
	}
	return nil
}
