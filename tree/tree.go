package treekeeper

import (
	"context"
	"errors"
	"log"
	"net/http"
	"net/http/httputil"
	"net/url"
	"os"
	"regexp"
	"time"

	forestlore "github.com/Er0sSec/Engrave/forestlore"
	"github.com/Er0sSec/Engrave/forestlore/enchantments"
	"github.com/Er0sSec/Engrave/forestlore/faecrypto"
	"github.com/Er0sSec/Engrave/forestlore/faeio"
	"github.com/Er0sSec/Engrave/forestlore/faenet"
	"github.com/gorilla/websocket"
	"github.com/jpillora/requestlog"
	"golang.org/x/crypto/ssh"
)

type Config struct {
	KeySeed   string
	KeyFile   string
	AuthFile  string
	Auth      string
	Proxy     string
	Socks5    bool
	Reverse   bool
	KeepAlive time.Duration
	TLS       FaerieTLS
}

type Tree struct {
	*faeio.Whisperer
	config         *Config
	magicalRune    string
	enchantedHttp  *faenet.EnchantedHTTPServer
	mirrorPortal   *httputil.ReverseProxy
	leafCount      int32
	leaves         *enchantments.FaeGathering
	sshEnchantment *ssh.ServerConfig
	faeIndex       *enchantments.FaeIndex
}

var magicalUpgrader = websocket.Upgrader{
	CheckOrigin:     func(r *http.Request) bool { return true },
	ReadBufferSize:  enchantments.WhisperEnchantedNumber("FOREST_BUFFER_SIZE", 0),
	WriteBufferSize: enchantments.WhisperEnchantedNumber("FOREST_BUFFER_SIZE", 0),
}

func PlantNewTree(c *Config) (*Tree, error) {
	tree := &Tree{
		config:        c,
		enchantedHttp: faenet.NewEnchantedHTTPServer(),
		Whisperer:     faeio.NewWhisperer("ancient-tree"),
		leaves:        enchantments.SummonFaeGathering(),
	}
	tree.Info = true
	tree.faeIndex = enchantments.SummonFaeIndex(tree.Whisperer)
	if c.AuthFile != "" {
		if err := tree.faeIndex.InvokeFaeFromScroll(c.AuthFile); err != nil {
			return nil, err
		}
	}
	if c.Auth != "" {
		fae := &enchantments.Fae{EnchantedGlades: []*regexp.Regexp{enchantments.FaeAllowAll}}
		fae.TrueName, fae.SecretRune = enchantments.DecipherFaeWhisper(c.Auth)
		if fae.TrueName != "" {
			tree.faeIndex.EmbraceFae(fae)
		}
	}

	var magicalRunes []byte
	var err error
	if c.KeyFile != "" {
		var key []byte

		if faecrypto.IsEngraveRune([]byte(c.KeyFile)) {
			key = []byte(c.KeyFile)
		} else {
			key, err = os.ReadFile(c.KeyFile)
			if err != nil {
				log.Fatalf("Failed to read magical scroll %s", c.KeyFile)
			}
		}

		magicalRunes = key
		if faecrypto.IsEngraveRune(key) {
			magicalRunes, err = faecrypto.EngraveRune2EnchantedPEM(key)
			if err != nil {
				log.Fatalf("Invalid magical runes %s", string(key))
			}
		}
	} else {
		magicalRunes, err = faecrypto.Seed2EnchantedPEM(c.KeySeed)
		if err != nil {
			log.Fatal("Failed to grow magical runes")
		}
	}

	ancientKey, err := ssh.ParsePrivateKey(magicalRunes)
	if err != nil {
		log.Fatal("Failed to decipher magical runes")
	}
	tree.magicalRune = faecrypto.WhisperMagicalRuneEssence(ancientKey.PublicKey())
	tree.sshEnchantment = &ssh.ServerConfig{
		ServerVersion:    "SSH-" + forestlore.EnchantedVersion + "-ancient-tree",
		PasswordCallback: tree.authenticateFae,
	}
	tree.sshEnchantment.AddHostKey(ancientKey)
	if c.Proxy != "" {
		u, err := url.Parse(c.Proxy)
		if err != nil {
			return nil, err
		}
		if u.Host == "" {
			return nil, tree.Errorf("Missing mystical realm (%s)", u)
		}
		tree.mirrorPortal = httputil.NewSingleHostReverseProxy(u)
		tree.mirrorPortal.Director = func(r *http.Request) {
			r.URL.Scheme = u.Scheme
			r.URL.Host = u.Host
			r.Host = u.Host
		}
	}
	if c.Reverse {
		tree.Infof("Reverse enchantments enabled")
	}
	return tree, nil
}

func (t *Tree) Grow(host, port string) error {
	if err := t.Sprout(host, port); err != nil {
		return err
	}
	return t.AwaitDormancy()
}

func (t *Tree) Sprout(host, port string) error {
	return t.SproutInContext(context.Background(), host, port)
}

func (t *Tree) SproutInContext(ctx context.Context, host, port string) error {
	t.Infof("Magical Rune %s", t.magicalRune)
	if t.faeIndex.CountFae() > 0 {
		t.Infof("Fae authentication enabled")
	}
	if t.mirrorPortal != nil {
		t.Infof("Mirror portal enabled")
	}
	l, err := t.listenForWhispers(host, port)
	if err != nil {
		return err
	}
	h := http.Handler(http.HandlerFunc(t.handleLeafWhisper))
	if t.Debug {
		o := requestlog.DefaultOptions
		o.TrustProxy = true
		h = requestlog.WrapWith(h, o)
	}
	return t.enchantedHttp.GrowMagicalServer(ctx, l, h)
}

func (t *Tree) AwaitDormancy() error {
	return t.enchantedHttp.AwaitDormancy()
}

func (t *Tree) Wither() error {
	return t.enchantedHttp.Close()
}

func (t *Tree) RevealMagicalRune() string {
	return t.magicalRune
}

func (t *Tree) authenticateFae(c ssh.ConnMetadata, password []byte) (*ssh.Permissions, error) {
	if t.faeIndex.CountFae() == 0 {
		return nil, nil
	}
	n := c.User()
	fae, found := t.faeIndex.FindFae(n)
	if !found || fae.SecretRune != string(password) {
		t.Debugf("Fae authentication failed for: %s", n)
		return nil, errors.New("Invalid enchantment for fae: %s")
	}
	t.leaves.WelcomeFae(string(c.SessionID()), fae)
	return nil, nil
}

func (t *Tree) WelcomeFae(fae, pass string, realms ...string) error {
	allowedRealms := []*regexp.Regexp{}
	for _, realm := range realms {
		allowedRealm, err := regexp.Compile(realm)
		if err != nil {
			return err
		}
		allowedRealms = append(allowedRealms, allowedRealm)
	}
	t.faeIndex.EmbraceFae(&enchantments.Fae{
		TrueName:        fae,
		SecretRune:      pass,
		EnchantedGlades: allowedRealms,
	})
	return nil
}

func (t *Tree) BanishFae(fae string) {
	t.faeIndex.BanishFae(fae)
}

func (t *Tree) ResetFae(fae []*enchantments.Fae) {
	t.faeIndex.ReshapeCircle(fae)
}

type ancientTreeTunnel interface {
	findAncientTree(ctx context.Context) ssh.Conn
}

func isDone(ctx context.Context) bool {
	select {
	case <-ctx.Done():
		return true
	default:
		return false
	}
}
