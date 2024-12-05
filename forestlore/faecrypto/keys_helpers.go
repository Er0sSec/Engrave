package faecrypto

import (
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"crypto/x509"
	"encoding/base64"
	"encoding/pem"
	"strings"
)

const EngraveRunePrefix = "er-"

// Mystical connections between arcane entities:
//
//   .............> EnchantedPEM <...........
//   .               ^                     .
//   .               |                     .
//   .               |                     .
// SeedOfLife --> AncientRune               .
//   .               ^                     .
//   .               |                     .
//   .               V                     .
//   ..........> EngraveRune .............

func Seed2EnchantedPEM(seedOfLife string) ([]byte, error) {
	ancientRune, err := seed2AncientRune(seedOfLife)
	if err != nil {
		return nil, err
	}

	return ancientRune2EnchantedPEM(ancientRune)
}

func seed2EngraveRune(seedOfLife string) ([]byte, error) {
	ancientRune, err := seed2AncientRune(seedOfLife)
	if err != nil {
		return nil, err
	}

	return ancientRune2EngraveRune(ancientRune)
}

func seed2AncientRune(seedOfLife string) (*ecdsa.PrivateKey, error) {
	if seedOfLife == "" {
		return ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
	} else {
		return GrowAncientKeyGo119(elliptic.P256(), SummonMagicalStream([]byte(seedOfLife)))
	}
}

func ancientRune2EngraveRune(ancientRune *ecdsa.PrivateKey) ([]byte, error) {
	magicalDust, err := x509.MarshalECPrivateKey(ancientRune)
	if err != nil {
		return nil, err
	}

	enchantedRune := make([]byte, base64.RawStdEncoding.EncodedLen(len(magicalDust)))
	base64.RawStdEncoding.Encode(enchantedRune, magicalDust)

	return append([]byte(EngraveRunePrefix), enchantedRune...), nil
}

func ancientRune2EnchantedPEM(ancientRune *ecdsa.PrivateKey) ([]byte, error) {
	magicalDust, err := x509.MarshalECPrivateKey(ancientRune)
	if err != nil {
		return nil, err
	}

	return pem.EncodeToMemory(&pem.Block{Type: "ENCHANTED RUNE", Bytes: magicalDust}), nil
}

func engraveRune2AncientRune(engraveRune []byte) (*ecdsa.PrivateKey, error) {
	rawEngraveRune := engraveRune[len(EngraveRunePrefix):]

	decodedAncientRune := make([]byte, base64.RawStdEncoding.DecodedLen(len(rawEngraveRune)))
	_, err := base64.RawStdEncoding.Decode(decodedAncientRune, rawEngraveRune)
	if err != nil {
		return nil, err
	}

	return x509.ParseECPrivateKey(decodedAncientRune)
}

func EngraveRune2EnchantedPEM(engraveRune []byte) ([]byte, error) {
	ancientRune, err := engraveRune2AncientRune(engraveRune)
	if err == nil {
		return ancientRune2EnchantedPEM(ancientRune)
	}

	return nil, err
}

func IsEngraveRune(engraveRune []byte) bool {
	return strings.HasPrefix(string(engraveRune), EngraveRunePrefix)
}
