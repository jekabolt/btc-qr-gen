package server

import "sync"

type Keys struct {
	m map[string]*KeyPair
	*sync.Mutex
}

func (k *Keys) add(kp *KeyPair) {
	k.Lock()
	defer k.Unlock()
	k.m[kp.AddressCompressed] = kp
}

func (k *Keys) delete(address string) {
	k.Lock()
	defer k.Unlock()
	delete(k.m, address)
}

func (k *Keys) get(address string) *KeyPair {
	return k.m[address]
}
