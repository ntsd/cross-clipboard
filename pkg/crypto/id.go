package crypto

import (
	"encoding/pem"
	"fmt"

	"github.com/libp2p/go-libp2p-core/crypto"
)

const (
	keySize int = 1024
	keyType int = crypto.Ed25519
)

// NewKeyPair generate p2p id key pair
func NewKeyPair() (crypto.PrivKey, crypto.PubKey, error) {
	return crypto.GenerateKeyPair(keyType, keySize)
}

// MarshalIDPrivateKey will return the p2p id private key in pem encoded
func MarshalIDPrivateKey(prvKey crypto.PrivKey) (string, error) {
	x509Encoded, err := crypto.MarshalPrivateKey(prvKey)
	if err != nil {
		return "", fmt.Errorf("unable to marshal id private key: %w", err)
	}
	pemEncoded := pem.EncodeToMemory(&pem.Block{Type: "ID PRIVATE KEY", Bytes: x509Encoded})
	return string(pemEncoded), nil
}

// UnmarshalIDPrivateKey receive p2p id private key in pem encoded
func UnmarshalIDPrivateKey(pemEncoded string) (crypto.PrivKey, error) {
	block, _ := pem.Decode([]byte(pemEncoded))
	x509Encoded := block.Bytes
	return crypto.UnmarshalPrivateKey(x509Encoded)
}

func GenerateIDPem() (string, error) {
	prvKey, _, err := NewKeyPair()
	if err != nil {
		return "", err
	}
	pkPem, err := MarshalIDPrivateKey(prvKey)
	if err != nil {
		return "", err
	}
	return pkPem, nil
}
