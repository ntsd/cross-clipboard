package crypto

import (
	"reflect"
	"testing"
)

func TestPGP(t *testing.T) {
	pem, err := GeneratePGPKey("test")
	if err != nil {
		t.Fatal(err)
	}

	pgpPrivKey, err := UnmarshalPGPKey(pem, nil)
	if err != nil {
		t.Fatal(err)
	}

	pubKey, err := pgpPrivKey.GetPublicKey()
	if err != nil {
		t.Fatal(err)
	}
	pgpPubKey, err := ByteToPGPKey(pubKey)
	if err != nil {
		t.Fatal(err)
	}

	pgpDecrypter, err := NewPGPDecrypter(pgpPrivKey)
	if err != nil {
		t.Fatal(err)
	}
	pgpEncrypter, err := NewPGPEncrypter(pgpPubKey)
	if err != nil {
		t.Fatal(err)
	}

	tests := []struct {
		message []byte
	}{
		{
			message: []byte("test secret"),
		},
	}

	for _, test := range tests {

		encrypted, err := pgpEncrypter.EncryptMessage(test.message)
		if err != nil {
			t.Fatal(err)
		}

		got, err := pgpDecrypter.DecryptMessage(encrypted)
		if err != nil {
			t.Fatal(err)
		}

		if !reflect.DeepEqual(got, test.message) {
			t.Errorf("got %q, wanted %q", got, test.message)
		}
	}
}
