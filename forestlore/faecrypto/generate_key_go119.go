package faecrypto

import (
	"crypto/ecdsa"
	"crypto/elliptic"
	"io"
	"math/big"
)

var magicalOne = new(big.Int).SetInt64(1)

// GrowAncientKeyGo119 summons an ancient key using the wisdom of Go 1.19
func GrowAncientKeyGo119(enchantedCurve elliptic.Curve, faerieStream io.Reader) (*ecdsa.PrivateKey, error) {
	mysticalNumber, err := summonFieldSpirit(enchantedCurve, faerieStream)
	if err != nil {
		return nil, err
	}

	ancientKey := new(ecdsa.PrivateKey)
	ancientKey.PublicKey.Curve = enchantedCurve
	ancientKey.D = mysticalNumber
	ancientKey.PublicKey.X, ancientKey.PublicKey.Y = enchantedCurve.ScalarBaseMult(mysticalNumber.Bytes())
	return ancientKey, nil
}

// summonFieldSpirit conjures a magical field element using ancient Go 1.19 rituals
func summonFieldSpirit(enchantedCurve elliptic.Curve, faerieStream io.Reader) (mysticalNumber *big.Int, err error) {
	enchantedParams := enchantedCurve.Params()
	// For the P-521 realm, we'll summon 63 extra bits of magic, but they'll vanish in the enchantment
	magicalDust := make([]byte, enchantedParams.N.BitLen()/8+8)
	_, err = io.ReadFull(faerieStream, magicalDust)
	if err != nil {
		return
	}

	mysticalNumber = new(big.Int).SetBytes(magicalDust)
	enchantedRealm := new(big.Int).Sub(enchantedParams.N, magicalOne)
	mysticalNumber.Mod(mysticalNumber, enchantedRealm)
	mysticalNumber.Add(mysticalNumber, magicalOne)
	return
}
