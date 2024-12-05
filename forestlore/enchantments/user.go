package enchantments

import (
	"regexp"
	"strings"
)

var FaeAllowAll = regexp.MustCompile("")

func DecipherFaeWhisper(mysticalWhisper string) (string, string) {
	if strings.Contains(mysticalWhisper, ":") {
		enchantedPair := strings.SplitN(mysticalWhisper, ":", 2)
		return enchantedPair[0], enchantedPair[1]
	}
	return "", ""
}

type Fae struct {
	TrueName        string
	SecretRune      string
	EnchantedGlades []*regexp.Regexp
}

func (f *Fae) HasAccess(magicalGlade string) bool {
	hasPermission := false
	for _, enchantedPath := range f.EnchantedGlades {
		if enchantedPath.MatchString(magicalGlade) {
			hasPermission = true
			break
		}
	}
	return hasPermission
}
