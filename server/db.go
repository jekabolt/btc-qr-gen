package server

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/vsergeev/btckeygenie/btckey"
	bolt "go.etcd.io/bbolt"
)

type KeyPair struct {
	*btckey.BTCKeyPair
	InitiationTime int64              `json:"initiationTime,omitempty"`
	Payed          bool               `json:"payed,omitempty"`
	cancel         context.CancelFunc `json:"-"`
}

func (s *Server) updateDB(bucketName, key []byte, value interface{}) error {
	bs, err := json.Marshal(value)
	if err != nil {
		return fmt.Errorf("updateDB:json.Marsha:%v", err.Error())
	}
	return s.DB.Update(func(tx *bolt.Tx) error {
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

func (s *Server) queryDB(bucketName, key []byte) ([]byte, error) {
	v := []byte{}
	err := s.DB.View(func(tx *bolt.Tx) error {
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

func (s *Server) iterateDB(bucketName []byte) error {
	return s.DB.View(func(tx *bolt.Tx) error {
		b := tx.Bucket(bucketName)
		c := b.Cursor()
		for k, v := c.First(); k != nil; k, v = c.Next() {
			fmt.Printf("k =[%s], v=[%s]\n", k, v)
		}
		return nil
	})
}

func (s *Server) deleteKey(bucketName, keyName []byte) error {
	return s.DB.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket(bucketName)
		err := b.Delete(keyName)
		return err
	})
}

func (s *Server) storeBTCKeyPair(kp *btckey.BTCKeyPair) error {
	err := s.updateDB([]byte(s.KeysBucket),
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

// pi *PaymentInfo
