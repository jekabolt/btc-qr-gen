package server

import (
	"encoding/binary"
	"encoding/json"
	"fmt"
	"time"

	"github.com/vsergeev/btckeygenie/btckey"
	bolt "go.etcd.io/bbolt"
)

type KeyPair struct {
	*btckey.BTCKeyPair
	Id             int
	InitiationTime int64
	Payed          bool
}

func (s *Server) createBucket(bucketName string) error {
	return s.DB.Update(func(tx *bolt.Tx) error {
		_, err := tx.CreateBucket([]byte(bucketName))
		if err != nil {
			return fmt.Errorf("create bucket: %s", err)
		}
		return nil
	})
}

func (s *Server) addKeyPair(k *btckey.BTCKeyPair) error {
	return s.DB.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(s.KeysBucket))

		id, err := b.NextSequence()
		if err != nil {
			return err
		}

		kp := KeyPair{
			BTCKeyPair:     k,
			Id:             int(id),
			InitiationTime: time.Now().Unix(),
		}
		buf, err := json.Marshal(kp)
		if err != nil {
			return err
		}
		return b.Put(itob(kp.Id), buf)
	})
}

// itob returns an 8-byte big endian representation of v.
func itob(v int) []byte {
	b := make([]byte, 8)
	binary.BigEndian.PutUint64(b, uint64(v))
	return b
}
