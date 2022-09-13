package crypto

import (
	"github.com/ProtonMail/gopenpgp/v2/crypto"
)

const (
	email      string = ""
	pgpKeyType string = "x25519"
)

// GeneratePGPKey generate pgp key armored format
func GeneratePGPKey(name string) (string, error) {
	ecKey, err := crypto.GenerateKey(name, email, pgpKeyType, 0)
	if err != nil {
		return "", err
	}
	armor, err := ecKey.Armor()
	if err != nil {
		return "", err
	}
	return armor, nil
}

// UnmarshalPGPKey unmarshals armored pgp key
func UnmarshalPGPKey(armoredPrivkey string, passphrase *[]byte) (*crypto.Key, error) {
	pgpPrivateKey, err := crypto.NewKeyFromArmored(armoredPrivkey)
	if err != nil {
		return nil, err
	}

	if passphrase != nil {
		unlockedKeyObj, err := pgpPrivateKey.Unlock(*passphrase)
		if err != nil {
			return nil, err
		}
		pgpPrivateKey = unlockedKeyObj
	}

	return pgpPrivateKey, nil
}

// ByteToPGPKey unmarshals unarmored pgp key
func ByteToPGPKey(pubKey []byte) (*crypto.Key, error) {
	pgpPubKey, err := crypto.NewKey(pubKey)
	if err != nil {
		return nil, err
	}

	return pgpPubKey, nil
}
