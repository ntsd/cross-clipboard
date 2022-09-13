package crypto

import (
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"crypto/x509"
	"encoding/pem"
	"fmt"
)

var curve elliptic.Curve = elliptic.P256()

func GenerateECDSAPrivateKeyPem() (string, error) {
	privateKey, err := ecdsa.GenerateKey(curve, rand.Reader)
	if err != nil {
		return "", fmt.Errorf("error generating ecdsa private key: %w", err)
	}

	pkPem, err := MarshalECDSAPrivateKey(privateKey)
	if err != nil {
		return "", fmt.Errorf("error to unmarshal ecdsa private key: %w", err)
	}

	return pkPem, nil
}

// UnmarshalECDSAPrivateKey marshal ecdsa private key to pem encoded
func MarshalECDSAPrivateKey(privateKey *ecdsa.PrivateKey) (string, error) {
	x509Encoded, err := x509.MarshalECPrivateKey(privateKey)
	if err != nil {
		return "", fmt.Errorf("error to marshal ecdsa private key: %w", err)
	}
	pemEncoded := pem.EncodeToMemory(&pem.Block{Type: "PGP PRIVATE KEY", Bytes: x509Encoded})
	return string(pemEncoded), nil
}

// UnmarshalECDSAPrivateKey unmarshals ecdsa private key from pem encoded
func UnmarshalECDSAPrivateKey(pemEncoded string) (*ecdsa.PrivateKey, error) {
	block, _ := pem.Decode([]byte(pemEncoded))
	x509Encoded := block.Bytes
	return x509.ParseECPrivateKey(x509Encoded)
}
