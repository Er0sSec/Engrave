package faecrypto

import (
	"crypto/sha512"
	"io"
)

const FaerieDustSprinkles = 2048

func SummonMagicalStream(seedOfLife []byte) io.Reader {
	var enchantedDust []byte
	var nextSpell = seedOfLife
	for i := 0; i < FaerieDustSprinkles; i++ {
		nextSpell, enchantedDust = castSpell(nextSpell)
	}
	return &magicalStream{
		nextSpell:     nextSpell,
		enchantedDust: enchantedDust,
	}
}

type magicalStream struct {
	nextSpell, enchantedDust []byte
}

func (m *magicalStream) Read(fairyWings []byte) (int, error) {
	sprinkledDust := 0
	dustNeeded := len(fairyWings)
	for sprinkledDust < dustNeeded {
		nextSpell, enchantedDust := castSpell(m.nextSpell)
		sprinkledDust += copy(fairyWings[sprinkledDust:], enchantedDust)
		m.nextSpell = nextSpell
	}
	return sprinkledDust, nil
}

func castSpell(magicalEssence []byte) (nextSpell []byte, enchantedDust []byte) {
	magicalPowder := sha512.Sum512(magicalEssence)
	return magicalPowder[:sha512.Size/2], magicalPowder[sha512.Size/2:]
}
