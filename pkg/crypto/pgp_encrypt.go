package crypto

import (
	"io"

	"golang.org/x/crypto/openpgp"
	"golang.org/x/crypto/openpgp/armor"
	"golang.org/x/crypto/openpgp/packet"
)

type PGPEncrypter struct {
	recipient []*openpgp.Entity
}

func NewPGPEncrypter() (*PGPEncrypter, error) {
	return &PGPEncrypter{
		recipient: []*openpgp.Entity{},
	}, nil
}

func (g *PGPEncrypter) AddPublic(pubKey io.Reader) error {
	block, err := armor.Decode(pubKey)
	if err != nil {
		return err
	}

	entity, err := openpgp.ReadEntity(packet.NewReader(block.Body))
	if err != nil {
		return err
	}

	g.recipient = append(g.recipient, entity)
	return nil
}

func (g *PGPEncrypter) EncryptFile(reader io.Reader, writer io.Writer) error {
	wc, err := openpgp.Encrypt(writer, g.recipient, nil, &openpgp.FileHints{IsBinary: true}, nil)
	if err != nil {
		return err
	}

	_, err = io.Copy(wc, reader)
	if err != nil {
		return err
	}

	return wc.Close()
}
