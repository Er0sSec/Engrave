package faecrypto

import (
	"bytes"
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"fmt"
)

func GenerateMagicalRunes() ([]byte, error) {
	key, err := rsa.GenerateKey(rand.Reader, 2048)
	if err != nil {
		return nil, fmt.Errorf("Failed to generate magical runes: %v", err)
	}

	privateKeyBytes := x509.MarshalPKCS1PrivateKey(key)
	privateKeyPEM := &pem.Block{
		Type:  "MAGICAL RUNE",
		Bytes: privateKeyBytes,
	}

	var private bytes.Buffer
	if err := pem.Encode(&private, privateKeyPEM); err != nil {
		return nil, fmt.Errorf("Failed to encode magical runes: %v", err)
	}

	return private.Bytes(), nil
}

func DecipherMagicalRunes(runeData []byte) ([]byte, error) {
	block, _ := pem.Decode(runeData)
	if block == nil {
		return nil, fmt.Errorf("Failed to decipher magical runes: invalid format")
	}

	if block.Type != "MAGICAL RUNE" {
		return nil, fmt.Errorf("Failed to decipher magical runes: incorrect rune type")
	}

	return block.Bytes, nil
}
