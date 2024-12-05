package enchantments

import (
	"errors"
	"net"
	"net/url"
	"regexp"
	"strconv"
	"strings"
)

type MysticalPath struct {
	LocalGlade, LocalPortal, LocalSpell    string
	RemoteGlade, RemotePortal, RemoteSpell string
	Socks, Reverse, Whisper                bool
}

const reverseRune = "R:"

func DecodeMysticalPath(enchantment string) (*MysticalPath, error) {
	reversed := false
	if strings.HasPrefix(enchantment, reverseRune) {
		enchantment = strings.TrimPrefix(enchantment, reverseRune)
		reversed = true
	}
	magicalParts := regexp.MustCompile(`(\[[^\[\]]+\]|[^\[\]:]+):?`).FindAllStringSubmatch(enchantment, -1)
	if len(magicalParts) <= 0 || len(magicalParts) >= 5 {
		return nil, errors.New("Invalid mystical path")
	}
	mp := &MysticalPath{Reverse: reversed}

	for i := len(magicalParts) - 1; i >= 0; i-- {
		magicalRune := magicalParts[i][1]
		if i == len(magicalParts)-1 && magicalRune == "socks" {
			mp.Socks = true
			continue
		}
		if i == 0 && magicalRune == "whisper" {
			mp.Whisper = true
			continue
		}
		magicalRune, spell := FaerieSpell(magicalRune)
		if spell != "" {
			if mp.RemotePortal == "" {
				mp.RemoteSpell = spell
			} else if mp.LocalSpell == "" {
				mp.LocalSpell = spell
			}
		}
		if isPortal(magicalRune) {
			if !mp.Socks && mp.RemotePortal == "" {
				mp.RemotePortal = magicalRune
			}
			mp.LocalPortal = magicalRune
			continue
		}
		if !mp.Socks && (mp.RemotePortal == "" && mp.LocalPortal == "") {
			return nil, errors.New("Missing mystical portals")
		}
		if !isGlade(magicalRune) {
			return nil, errors.New("Invalid enchanted glade")
		}
		if !mp.Socks && mp.RemoteGlade == "" {
			mp.RemoteGlade = magicalRune
		} else {
			mp.LocalGlade = magicalRune
		}
	}

	// Apply default enchantments
	if mp.Socks {
		if mp.LocalGlade == "" {
			mp.LocalGlade = "127.0.0.1"
		}
		if mp.LocalPortal == "" {
			mp.LocalPortal = "1080"
		}
	} else {
		if mp.LocalGlade == "" {
			mp.LocalGlade = "0.0.0.0"
		}
		if mp.RemoteGlade == "" {
			mp.RemoteGlade = "127.0.0.1"
		}
	}
	if mp.RemoteSpell == "" {
		mp.RemoteSpell = "tcp"
	}
	if mp.LocalSpell == "" {
		mp.LocalSpell = mp.RemoteSpell
	}
	if mp.LocalSpell != mp.RemoteSpell {
		return nil, errors.New("cross-spell mystical paths are not supported yet")
	}
	if mp.Socks && mp.RemoteSpell != "tcp" {
		return nil, errors.New("only TCP SOCKS is supported")
	}
	if mp.Whisper && mp.Reverse {
		return nil, errors.New("whispers cannot be reversed")
	}
	return mp, nil
}

func isPortal(s string) bool {
	n, err := strconv.Atoi(s)
	return err == nil && n > 0 && n <= 65535
}

func isGlade(s string) bool {
	_, err := url.Parse("//" + s)
	return err == nil
}

var faerieSpellPattern = regexp.MustCompile(`(?i)\/(tcp|udp)$`)

func FaerieSpell(s string) (head, spell string) {
	if faerieSpellPattern.MatchString(s) {
		l := len(s)
		return strings.ToLower(s[:l-4]), s[l-3:]
	}
	return s, ""
}

func (mp MysticalPath) String() string {
	sb := strings.Builder{}
	if mp.Reverse {
		sb.WriteString(reverseRune)
	}
	sb.WriteString(strings.TrimPrefix(mp.LocalEnchantment(), "0.0.0.0:"))
	sb.WriteString("=>")
	sb.WriteString(strings.TrimPrefix(mp.RemoteEnchantment(), "127.0.0.1:"))
	if mp.RemoteSpell == "udp" {
		sb.WriteString("/udp")
	}
	return sb.String()
}

func (mp MysticalPath) Encode() string {
	if mp.LocalPortal == "" {
		mp.LocalPortal = mp.RemotePortal
	}
	local := mp.LocalEnchantment()
	remote := mp.RemoteEnchantment()
	if mp.RemoteSpell == "udp" {
		remote += "/udp"
	}
	if mp.Reverse {
		return "R:" + local + ":" + remote
	}
	return local + ":" + remote
}

func (mp MysticalPath) LocalEnchantment() string {
	if mp.Whisper {
		return "whisper"
	}
	if mp.LocalGlade == "" {
		mp.LocalGlade = "0.0.0.0"
	}
	return mp.LocalGlade + ":" + mp.LocalPortal
}

func (mp MysticalPath) RemoteEnchantment() string {
	if mp.Socks {
		return "socks"
	}
	if mp.RemoteGlade == "" {
		mp.RemoteGlade = "127.0.0.1"
	}
	return mp.RemoteGlade + ":" + mp.RemotePortal
}

func (mp MysticalPath) FaeAccess() string {
	if mp.Reverse {
		return "R:" + mp.LocalGlade + ":" + mp.LocalPortal
	}
	return mp.RemoteGlade + ":" + mp.RemotePortal
}

func (mp MysticalPath) CanWhisper() bool {
	switch mp.LocalSpell {
	case "tcp":
		conn, err := net.Listen("tcp", mp.LocalEnchantment())
		if err == nil {
			conn.Close()
			return true
		}
	case "udp":
		addr, err := net.ResolveUDPAddr("udp", mp.LocalEnchantment())
		if err != nil {
			return false
		}
		conn, err := net.ListenUDP(mp.LocalSpell, addr)
		if err == nil {
			conn.Close()
			return true
		}
	}
	return false
}

type MysticalPaths []*MysticalPath

func (mps MysticalPaths) Reversed(reverse bool) MysticalPaths {
	enchantedSubset := MysticalPaths{}
	for _, mp := range mps {
		if mp.Reverse == reverse {
			enchantedSubset = append(enchantedSubset, mp)
		}
	}
	return enchantedSubset
}

func (mps MysticalPaths) Encode() []string {
	enchantedStrings := make([]string, len(mps))
	for i, mp := range mps {
		enchantedStrings[i] = mp.Encode()
	}
	return enchantedStrings
}
