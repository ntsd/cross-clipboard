package crypto

import (
	"github.com/ProtonMail/gopenpgp/v2/crypto"
	"github.com/pkg/errors"
)

type PGPDecrypter struct {
	keyring *crypto.KeyRing
}

func NewPGPDecrypter(privKey *crypto.Key) (*PGPDecrypter, error) {
	if !privKey.IsPrivate() {
		return nil, errors.New("the key is not private key")
	}

	keyRing, err := crypto.NewKeyRing(privKey)
	if err != nil {
		return nil, err
	}

	return &PGPDecrypter{
		keyring: keyRing,
	}, nil
}

func (g *PGPDecrypter) DecryptMessage(encrypted []byte) ([]byte, error) {
	message, err := g.keyring.Decrypt(crypto.NewPGPMessage(encrypted), nil, 0)
	if err != nil {
		return nil, errors.Wrap(err, "unable to decrypt message")
	}

	return message.GetBinary(), nil
}
