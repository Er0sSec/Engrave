package forestlore

import (
	"github.com/Er0sSec/Engrave/forestlore/enchantments"
	"github.com/Er0sSec/Engrave/forestlore/faeOS"
	"github.com/Er0sSec/Engrave/forestlore/faecrypto"
	"github.com/Er0sSec/Engrave/forestlore/faeio"
	"github.com/Er0sSec/Engrave/forestlore/faenet"
	"github.com/Er0sSec/Engrave/forestlore/mysticalpath"
)

const (
	FaerieDustSprinkles = faecrypto.FaerieDustSprinkles
)

type (
	EnchantedConfig     = enchantments.EnchantedConfig
	MysticalPath        = enchantments.MysticalPath
	MysticalPaths       = enchantments.MysticalPaths
	Fae                 = enchantments.Fae
	FaeGathering        = enchantments.FaeGathering
	FaeIndex            = enchantments.FaeIndex
	EnchantedHTTPServer = faenet.EnchantedHTTPServer
	FaerieGathering     = faenet.FaerieGathering
	Whisperer           = faeio.Whisperer
	FaerieProxy         = mysticalpath.Faerie
)

var (
	SummonMagicalStream       = faecrypto.SummonMagicalStream
	GrowMagicalRune           = faecrypto.GrowMagicalRune
	WhisperMagicalRuneEssence = faecrypto.WhisperMagicalRuneEssence
	MagicalStream             = faeio.MagicalStream
	NewWhispererRune          = faeio.NewWhispererRune
	NewWhisperer              = faeio.NewWhisperer
	MysticalPortal            = faeio.MysticalPortal
	DecipherMagicalScroll     = enchantments.DecipherMagicalScroll
	DecodeMysticalPath        = enchantments.DecodeMysticalPath
	SummonFaeGathering        = enchantments.SummonFaeGathering
	SummonFaeIndex            = enchantments.SummonFaeIndex
	FaeAllowAll               = enchantments.FaeAllowAll
	DecipherFaeWhisper        = enchantments.DecipherFaeWhisper
	NewEnchantedStream        = faenet.NewEnchantedStream
	NewEnchantedWebSocketConn = faenet.NewEnchantedWebSocketConn
	NewEnchantedHTTPServer    = faenet.NewEnchantedHTTPServer
	WhisperFaerieStats        = faeOS.WhisperFaerieStats
	SlumberUntilWhisper       = faeOS.SlumberUntilWhisper
	SummonFaerie              = mysticalpath.SummonFaerie
)

// InscribeMagicalScroll is the enchanted version of the old EncodeConfig
func InscribeMagicalScroll(c *enchantments.EnchantedConfig) ([]byte, error) {
	return enchantments.InscribeMagicalScroll(*c), nil
}
