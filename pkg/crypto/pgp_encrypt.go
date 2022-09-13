package crypto

import (
	"github.com/ProtonMail/gopenpgp/v2/crypto"
	"github.com/pkg/errors"
)

type PGPEncrypter struct {
	keyring *crypto.KeyRing
}

func NewPGPEncrypter(pubKey *crypto.Key) (*PGPEncrypter, error) {
	if pubKey.IsPrivate() {
		return nil, errors.New("the key is not public key")
	}

	keyRing, err := crypto.NewKeyRing(pubKey)
	if err != nil {
		return nil, err
	}

	return &PGPEncrypter{
		keyring: keyRing,
	}, nil
}

func (g *PGPEncrypter) EncryptMessage(message []byte) ([]byte, error) {
	pgpMessage, err := g.keyring.Encrypt(crypto.NewPlainMessage(message), nil)
	if err != nil {
		return nil, errors.Wrap(err, "unable to encrypt message")
	}

	return pgpMessage.GetBinary(), nil
}
