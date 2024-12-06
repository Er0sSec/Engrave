package leafwhisper

import (
	"context"
	"crypto/md5"
	"crypto/tls"
	"crypto/x509"
	"encoding/base64"
	"errors"
	"fmt"
	"net"
	"net/http"
	"net/url"
	"os"
	"regexp"
	"strings"
	"time"

	forestlore "github.com/Er0sSec/Engrave/forestlore"
	"github.com/Er0sSec/Engrave/forestlore/enchantments"
	"github.com/Er0sSec/Engrave/forestlore/faecrypto"
	"github.com/Er0sSec/Engrave/forestlore/faeio"
	"github.com/Er0sSec/Engrave/forestlore/faenet"
	"github.com/Er0sSec/Engrave/forestlore/mysticalpath"
	"github.com/gorilla/websocket"

	"golang.org/x/crypto/ssh"
	"golang.org/x/net/proxy"
	"golang.org/x/sync/errgroup"
)

type LeafConfig struct {
	MagicalRune     string                                                            // was Fingerprint
	FaeWhisper      string                                                            // was Auth
	MagicalPulse    time.Duration                                                     // was KeepAlive
	MaxRevivalCount int                                                               // was MaxRetryCount
	MaxRevivalPause time.Duration                                                     // was MaxRetryInterval
	AncientTree     string                                                            // was Server
	MysticalPortal  string                                                            // was Proxy
	EnchantedPaths  []string                                                          // was Remotes
	MagicalSeals    http.Header                                                       // was Headers
	FaerieTLS       FaerieTLS                                                         // already themed
	WeaveConnection func(ctx context.Context, network, addr string) (net.Conn, error) // was DialContext
	EnhancedVision  bool                                                              // was Verbose
}
type FaerieTLS struct {
	SkipVerify bool
	CA         string
	Cert       string
	Key        string
	ServerName string
}

type Leaf struct {
	*faeio.Whisperer
	config   *LeafConfig
	computed struct {
		MagicalVersion string
		MysticalPaths  enchantments.MysticalPaths // This should be a slice type that has Reversed method
	}
	enchantedConfig *ssh.ClientConfig          // was sshConfig
	faerieShield    *tls.Config                // was tlsConfig
	portalURL       *url.URL                   // was proxyURL
	ancientTree     string                     // was server
	faerieCount     faenet.FaerieGathering     // already themed
	wither          func()                     // was stop
	faerieGroup     *errgroup.Group            // was eg
	enchantedPath   *mysticalpath.MysticalPath // already themed
}

func GrowNewLeaf(c *LeafConfig) (*Leaf, error) {
	if !strings.HasPrefix(c.AncientTree, "http") {
		c.AncientTree = "http://" + c.AncientTree
	}
	if c.MaxRevivalPause < time.Second {
		c.MaxRevivalPause = 5 * time.Minute
	}
	u, err := url.Parse(c.AncientTree)
	if err != nil {
		return nil, err
	}
	u.Scheme = strings.Replace(u.Scheme, "http", "ws", 1)
	if !regexp.MustCompile(`:\d+$`).MatchString(u.Host) {
		if u.Scheme == "wss" {
			u.Host = u.Host + ":443"
		} else {
			u.Host = u.Host + ":80"
		}
	}
	hasReverse := false
	hasSocks := false
	hasStdio := false
	leaf := &Leaf{Whisperer: faeio.NewWhisperer("leaf"), config: c, computed: enchantments.EnchantedConfig{
		MagicalVersion: forestlore.EnchantedVersion,
	}, ancientTree: u.String(), faerieShield: nil}
	leaf.Whisperer.Info = true

	if u.Scheme == "wss" {
		tc := &tls.Config{}
		if c.FaerieTLS.ServerName != "" {
			tc.ServerName = c.FaerieTLS.ServerName
		}
		if c.FaerieTLS.SkipVerify {
			leaf.Infof("🧚 TLS verification disabled")
			tc.InsecureSkipVerify = true
		} else if c.FaerieTLS.CA != "" {
			rootCAs := x509.NewCertPool()
			if b, err := os.ReadFile(c.FaerieTLS.CA); err != nil {
				return nil, fmt.Errorf("🍄 Failed to load magical scroll: %s", c.FaerieTLS.CA)
			} else if ok := rootCAs.AppendCertsFromPEM(b); !ok {
				return nil, fmt.Errorf("🍄 Failed to decode magical runes: %s", c.FaerieTLS.CA)
			} else {
				leaf.Infof("🧚 TLS verification using magical scroll %s", c.FaerieTLS.CA)
				tc.RootCAs = rootCAs
			}
		}
		if c.FaerieTLS.Cert != "" && c.FaerieTLS.Key != "" {
			c, err := tls.LoadX509KeyPair(c.FaerieTLS.Cert, c.FaerieTLS.Key)
			if err != nil {
				return nil, fmt.Errorf("🍄 Error loading leaf's magical runes: %v", err)
			}
			tc.Certificates = []tls.Certificate{c}
		} else if c.FaerieTLS.Cert != "" || c.FaerieTLS.Key != "" {
			return nil, fmt.Errorf("🍄 Please provide BOTH magical runes for the leaf")
		}
		leaf.faerieShield = tc
	}

	for _, s := range c.EnchantedPaths {
		r, err := enchantments.DecodeRemote(s)
		if err != nil {
			return nil, fmt.Errorf("🍄 Failed to decode mystical pathway '%s': %s", s, err)
		}
		if r.Socks {
			hasSocks = true
		}
		if r.Reverse {
			hasReverse = true
		}
		if r.Whisper {
			if hasStdio {
				return nil, errors.New("🍄 Only one mystical stream is allowed")
			}
			hasStdio = true
		}
		if !r.Reverse && !r.Whisper && !r.CanWhisper() {
			return nil, fmt.Errorf("🍄 Leaf cannot listen on %s", r.String())
		}
		leaf.computed.MysticalPaths = append(leaf.computed.MysticalPaths, r)

	}

	if p := c.MysticalPortal; p != "" {
		leaf.portalURL, err = url.Parse(p)
		if err != nil {
			return nil, fmt.Errorf("🍄 Invalid mystical portal URL (%s)", err)
		}
	}

	user, pass := enchantments.DecipherFaeWhisper(c.FaeWhisper)
	leaf.enchantedConfig = &ssh.ClientConfig{
		User:            user,
		Auth:            []ssh.AuthMethod{ssh.Password(pass)},
		ClientVersion:   "SSH-" + forestlore.EnchantedVersion + "-leaf",
		HostKeyCallback: leaf.verifyTree,
		Timeout:         enchantments.WhisperTimespell("SSH_TIMEOUT", 30*time.Second),
	}

	leaf.enchantedPath = mysticalpath.New(mysticalpath.EnchantedConfig{
		Whisperer:     leaf.Whisperer,
		InboundMagic:  true,
		OutboundMagic: hasReverse,
		FaerieSocks:   hasReverse && hasSocks,
		MagicalPulse:  leaf.config.MagicalPulse,
	})
	return leaf, nil
}

func (l *Leaf) Sprout(context.Context) error {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	if err := l.GrowLeaves(ctx); err != nil {
		return err
	}
	return l.AwaitDormancy()
}

func (l *Leaf) verifyTree(hostname string, remote net.Addr, key ssh.PublicKey) error {
	expect := l.config.MagicalRune
	if expect == "" {
		return nil
	}
	got := faecrypto.WhisperMagicalRuneEssence(key)
	_, err := base64.StdEncoding.DecodeString(expect)
	if _, ok := err.(base64.CorruptInputError); ok {
		l.Whisperer.Infof("🍄 Specified outdated MD5 rune (%s), please update to the new SHA256 rune: %s", expect, got)
		return l.verifyAncientRune(key)
	} else if err != nil {
		return fmt.Errorf("🍄 Error decoding magical rune: %w", err)
	}
	if got != expect {
		return fmt.Errorf("🍄 Invalid magical rune (%s)", got)
	}
	l.Infof("🧚 Magical rune %s", got)
	return nil
}

func (l *Leaf) verifyAncientRune(key ssh.PublicKey) error {
	bytes := md5.Sum(key.Marshal())
	strbytes := make([]string, len(bytes))
	for i, b := range bytes {
		strbytes[i] = fmt.Sprintf("%02x", b)
	}
	got := strings.Join(strbytes, ":")
	expect := l.config.MagicalRune
	if !strings.HasPrefix(got, expect) {
		return fmt.Errorf("🍄 Invalid magical rune (%s)", got)
	}
	return nil
}

func (l *Leaf) GrowLeaves(ctx context.Context) error {
	ctx, cancel := context.WithCancel(ctx)
	l.wither = cancel
	eg, ctx := errgroup.WithContext(ctx)
	l.faerieGroup = eg
	via := ""
	if l.portalURL != nil {
		via = " via " + l.portalURL.String()
	}
	l.Infof("🌿 Connecting to %s%s\n", l.ancientTree, via)
	eg.Go(func() error {
		return l.magicalConnectionDance(ctx)
	})
	eg.Go(func() error {
		leafInbound := l.computed.MysticalPaths.Reversed(false)
		if len(leafInbound) == 0 {
			return nil
		}
		return l.enchantedPath.BindRemotes(ctx, leafInbound)
	})
	return nil
}

func (l *Leaf) setMysticalPortal(u *url.URL, d *websocket.Dialer) error {
	if !strings.HasPrefix(u.Scheme, "socks") {
		d.Proxy = func(*http.Request) (*url.URL, error) {
			return u, nil
		}
		return nil
	}
	if u.Scheme != "socks" && u.Scheme != "socks5h" {
		return fmt.Errorf(
			"🍄 unsupported socks mystical portal type: %s:// (only socks5h:// or socks:// is supported)",
			u.Scheme,
		)
	}
	var auth *proxy.Auth
	if u.User != nil {
		pass, _ := u.User.Password()
		auth = &proxy.Auth{
			User:     u.User.Username(),
			Password: pass,
		}
	}
	socksDialer, err := proxy.SOCKS5("tcp", u.Host, auth, proxy.Direct)
	if err != nil {
		return err
	}
	d.NetDial = socksDialer.Dial
	return nil
}

func (l *Leaf) AwaitDormancy() error {
	return l.faerieGroup.Wait()
}

func (l *Leaf) Wither() error {
	if l.wither != nil {
		l.wither()
	}
	return nil
}
