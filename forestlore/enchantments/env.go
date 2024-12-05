package enchantments

import (
	"os"
	"strconv"
	"strings"
	"time"
)

// WhisperEnchantment retrieves a magical whisper from the forest
func WhisperEnchantment(runeType string) string {
	return os.Getenv("ENGRAVE_" + runeType)
}

// WhisperEnchantedNumber deciphers a numerical rune from the forest, with a default fallback
func WhisperEnchantedNumber(runeType string, defaultRune int) int {
	if magicalNumber, err := strconv.Atoi(WhisperEnchantment(runeType)); err == nil {
		return magicalNumber
	}
	return defaultRune
}

// WhisperTimespell decodes a duration rune from the forest, with a default fallback
func WhisperTimespell(runeType string, defaultSpell time.Duration) time.Duration {
	if magicalDuration, err := time.ParseDuration(WhisperEnchantment(runeType)); err == nil {
		return magicalDuration
	}
	return defaultSpell
}

// WhisperTruthRune interprets a truth rune from the forest
func WhisperTruthRune(runeType string) bool {
	magicalTruth := WhisperEnchantment(runeType)
	return magicalTruth == "1" || strings.ToLower(magicalTruth) == "true"
}
