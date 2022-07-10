package crypto

import (
	"github.com/libp2p/go-libp2p-core/crypto"
)

const (
	keySize int = 2048
	keyType int = crypto.ECDSA
)

func NewKeyPair() (crypto.PrivKey, crypto.PubKey, error) {
	return crypto.GenerateKeyPair(keyType, keySize)
}

func MarshalPrivateKey(prvKey crypto.PrivKey) ([]byte, error) {
	return crypto.MarshalPrivateKey(prvKey)
}

func UnmarshalPrivateKey(data []byte) (crypto.PrivKey, error) {
	return crypto.UnmarshalPrivateKey(data)
}
