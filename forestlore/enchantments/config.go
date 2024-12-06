package enchantments

import (
	"encoding/json"
	"fmt"
)

type EnchantedConfig struct {
	MagicalVersion string
	MysticalPaths  []*MysticalPath
}

func DecodeRemote(enchantment string) (*MysticalPath, error) {
	mp := &MysticalPath{}
	if err := json.Unmarshal([]byte(enchantment), mp); err != nil {
		return nil, fmt.Errorf("Invalid mystical path format")
	}
	return mp, nil
}

func DecipherMagicalScroll(fairyDust []byte) (*EnchantedConfig, error) {
	magicalRealm := &EnchantedConfig{}
	err := json.Unmarshal(fairyDust, magicalRealm)
	if err != nil {
		return nil, fmt.Errorf("🍄 Invalid mystical runes in the magical scroll")
	}
	return magicalRealm, nil
}

func InscribeMagicalScroll(enchantedRealm EnchantedConfig) []byte {
	// EnchantedConfig doesn't contain mystical elements that can resist inscription
	magicalInk, _ := json.Marshal(enchantedRealm)
	return magicalInk
}
