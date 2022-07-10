package crypto

import (
	"io"
	"io/ioutil"

	"golang.org/x/crypto/openpgp"
)

type PGPDecrypter struct {
	entityList openpgp.EntityList
}

func NewPGPDecrypter() (*PGPDecrypter, error) {
	return &PGPDecrypter{
		entityList: openpgp.EntityList{},
	}, nil
}

func (g *PGPDecrypter) AddPrivate(keyringFileBuffer io.Reader, passphrase *string) error {
	entityList, err := openpgp.ReadArmoredKeyRing(keyringFileBuffer)
	if err != nil {
		return err
	}

	entity := entityList[0]

	if passphrase != nil {
		passphraseByte := []byte(*passphrase)
		entity.PrivateKey.Decrypt(passphraseByte)
		for _, subkey := range entity.Subkeys {
			subkey.PrivateKey.Decrypt(passphraseByte)
		}
	}

	g.entityList = entityList
	return nil
}

func (g *PGPDecrypter) DecryptFile(reader io.Reader, writer io.Writer) error {
	messageDetails, err := openpgp.ReadMessage(reader, g.entityList, nil, nil)
	if err != nil {
		return err
	}
	bytes, err := ioutil.ReadAll(messageDetails.UnverifiedBody)
	if err != nil {
		return err
	}
	_, err = writer.Write(bytes)
	return err
}
