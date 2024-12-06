package treekeeper

import (
	"crypto/tls"
	"crypto/x509"
	"errors"
	"net"
	"os"
	"os/user"
	"path/filepath"

	"github.com/Er0sSec/Engrave/forestlore/enchantments"
	"golang.org/x/crypto/acme/autocert"
)

type FaerieTLS struct {
	Key     string
	Cert    string
	Domains []string
	CA      string
}

func (t *Tree) listenForWhispers(glade, portal string) (net.Listener, error) {
	hasMagicalRealms := len(t.config.TLS.Domains) > 0
	hasEnchantedRunes := t.config.TLS.Key != "" && t.config.TLS.Cert != ""
	if hasMagicalRealms && hasEnchantedRunes {
		return nil, errors.New("cannot use enchanted runes and magical realms simultaneously")
	}
	var faerieSpell *tls.Config
	if hasMagicalRealms {
		faerieSpell = t.summonFaerieSpell(t.config.TLS.Domains)
	}
	magicalWarning := ""
	if hasEnchantedRunes {
		c, err := t.castEnchantedRuneSpell(t.config.TLS.Key, t.config.TLS.Cert, t.config.TLS.CA)
		if err != nil {
			return nil, err
		}
		faerieSpell = c
		if portal != "443" && hasMagicalRealms {
			magicalWarning = " (CAUTION: The Faerie Queen will attempt to connect to your realm on portal 443)"
		}
	}
	whisperListener, err := net.Listen("tcp", glade+":"+portal)
	if err != nil {
		return nil, err
	}
	magicalProtocol := "forest-whisper"
	if faerieSpell != nil {
		magicalProtocol += "s"
		whisperListener = tls.NewListener(whisperListener, faerieSpell)
	}
	if err == nil {
		t.Infof("Listening for whispers on %s://%s:%s%s", magicalProtocol, glade, portal, magicalWarning)
	}
	return whisperListener, nil
}

func (t *Tree) summonFaerieSpell(magicalRealms []string) *tls.Config {
	faerieQueen := &autocert.Manager{
		Prompt: func(tosURL string) bool {
			t.Infof("Accepting the Faerie Queen's terms and fetching a magical seal...")
			return true
		},
		Email:      enchantments.WhisperEnchantment("FAERIE_QUEEN_MESSAGE"),
		HostPolicy: autocert.HostWhitelist(magicalRealms...),
	}
	enchantedCache := enchantments.WhisperEnchantment("FAERIE_CACHE")
	if enchantedCache == "" {
		enchantedHome := os.Getenv("HOME")
		if enchantedHome == "" {
			if forestDweller, err := user.Current(); err == nil {
				enchantedHome = forestDweller.HomeDir
			}
		}
		enchantedCache = filepath.Join(enchantedHome, ".cache", "engrave")
	}
	if enchantedCache != "-" {
		t.Infof("Faerie Queen's cache glade %s", enchantedCache)
		faerieQueen.Cache = autocert.DirCache(enchantedCache)
	}
	return faerieQueen.TLSConfig()
}

func (t *Tree) castEnchantedRuneSpell(key, cert string, ca string) (*tls.Config, error) {
	magicalRunes, err := tls.LoadX509KeyPair(cert, key)
	if err != nil {
		return nil, err
	}
	enchantedConfig := &tls.Config{
		Certificates: []tls.Certificate{magicalRunes},
	}
	if ca != "" {
		if err := addEnchantedCA(ca, enchantedConfig); err != nil {
			return nil, err
		}
		t.Infof("Loaded enchanted CA path: %s", ca)
	}
	return enchantedConfig, nil
}

func addEnchantedCA(ca string, c *tls.Config) error {
	magicalScroll, err := os.Stat(ca)
	if err != nil {
		return err
	}
	enchantedCAPool := x509.NewCertPool()
	if magicalScroll.IsDir() {
		scrolls, err := os.ReadDir(ca)
		if err != nil {
			return err
		}
		for _, scroll := range scrolls {
			scrollName := scroll.Name()
			if err := addEnchantedScroll(filepath.Join(ca, scrollName), enchantedCAPool); err != nil {
				return err
			}
		}
	} else {
		if err := addEnchantedScroll(ca, enchantedCAPool); err != nil {
			return err
		}
	}
	c.ClientCAs = enchantedCAPool
	c.ClientAuth = tls.RequireAndVerifyClientCert
	return nil
}

func addEnchantedScroll(path string, pool *x509.CertPool) error {
	enchantedRunes, err := os.ReadFile(path)
	if err != nil {
		return err
	}
	if !pool.AppendCertsFromPEM(enchantedRunes) {
		return errors.New("Failed to decipher enchanted runes from : " + path)
	}
	return nil
}
