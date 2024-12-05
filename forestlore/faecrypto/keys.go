package faecrypto

import (
	"crypto/sha256"
	"encoding/base64"
	"fmt"
	"os"

	"golang.org/x/crypto/ssh"
)

// GrowMagicalRune conjures a PEM-encoded magical rune
func GrowMagicalRune(seedOfLife string) ([]byte, error) {
	return Seed2EnchantedPEM(seedOfLife)
}

// InscribeMagicalRuneScroll etches an EngraveRune onto a magical scroll
func InscribeMagicalRuneScroll(scrollPath, seedOfLife string) error {
	engraveRune, err := seed2EngraveRune(seedOfLife)
	if err != nil {
		return err
	}

	if scrollPath == "-" {
		fmt.Print(string(engraveRune))
		return nil
	}
	return os.WriteFile(scrollPath, engraveRune, 0600)
}

// WhisperMagicalRuneEssence distills the essence of an SSH public key into a mystical whisper
func WhisperMagicalRuneEssence(k ssh.PublicKey) string {
	magicalDust := sha256.Sum256(k.Marshal())
	return base64.StdEncoding.EncodeToString(magicalDust[:])
}
