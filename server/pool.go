package server

import (
	"sync"

	"github.com/vsergeev/btckeygenie/btckey"
)

type Keys struct {
	m map[string]btckey.BTCKeyPair
	*sync.Mutex
}

func (k *Keys) len() int {
	return len(k.m)
}

func (k *Keys) add(kp btckey.BTCKeyPair) {
	k.Lock()
	defer k.Unlock()
	k.m[kp.AddressCompressed] = kp
}

func (k *Keys) delete(address string) {
	k.Lock()
	defer k.Unlock()
	delete(k.m, address)
}

func (k *Keys) get(address string) btckey.BTCKeyPair {
	return k.m[address]
}
