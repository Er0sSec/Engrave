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

type Config struct {
	Fingerprint      string
	Auth             string
	KeepAlive        time.Duration
	MaxRetryCount    int
	MaxRetryInterval time.Duration
	Server           string
	Proxy            string
	Remotes          []string
	Headers          http.Header
	TLS              FaerieTLS
	DialContext      func(ctx context.Context, network, addr string) (net.Conn, error)
	Verbose          bool
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
	config       *Config
	computed     enchantments.Config
	sshConfig    *ssh.ClientConfig
	tlsConfig    *tls.Config
	proxyURL     *url.URL
	server       string
	connCount    faenet.ConnCount
	stop         func()
	eg           *errgroup.Group
	mysticalpath *mysticalpath.MysticalPath
}

func GrowNewLeaf(c *Config) (*Leaf, error) {
	if !strings.HasPrefix(c.Server, "http") {
		c.Server = "http://" + c.Server
	}
	if c.MaxRetryInterval < time.Second {
		c.MaxRetryInterval = 5 * time.Minute
	}
	u, err := url.Parse(c.Server)
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
	leaf := &Leaf{
		Whisperer: faeio.NewWhisperer("leaf"),
		config:    c,
		computed: enchantments.Config{
			Version: forestlore.EnchantedVersion,
		},
		server:    u.String(),
		tlsConfig: nil,
	}
	leaf.Whisperer.Info = true

	if u.Scheme == "wss" {
		tc := &tls.Config{}
		if c.TLS.ServerName != "" {
			tc.ServerName = c.TLS.ServerName
		}
		if c.TLS.SkipVerify {
			leaf.Infof("🧚 TLS verification disabled")
			tc.InsecureSkipVerify = true
		} else if c.TLS.CA != "" {
			rootCAs := x509.NewCertPool()
			if b, err := os.ReadFile(c.TLS.CA); err != nil {
				return nil, fmt.Errorf("🍄 Failed to load magical scroll: %s", c.TLS.CA)
			} else if ok := rootCAs.AppendCertsFromPEM(b); !ok {
				return nil, fmt.Errorf("🍄 Failed to decode magical runes: %s", c.TLS.CA)
			} else {
				leaf.Infof("🧚 TLS verification using magical scroll %s", c.TLS.CA)
				tc.RootCAs = rootCAs
			}
		}
		if c.TLS.Cert != "" && c.TLS.Key != "" {
			c, err := tls.LoadX509KeyPair(c.TLS.Cert, c.TLS.Key)
			if err != nil {
				return nil, fmt.Errorf("🍄 Error loading leaf's magical runes: %v", err)
			}
			tc.Certificates = []tls.Certificate{c}
		} else if c.TLS.Cert != "" || c.TLS.Key != "" {
			return nil, fmt.Errorf("🍄 Please provide BOTH magical runes for the leaf")
		}
		leaf.tlsConfig = tc
	}

	for _, s := range c.Remotes {
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
		if r.Stdio {
			if hasStdio {
				return nil, errors.New("🍄 Only one mystical stream is allowed")
			}
			hasStdio = true
		}
		if !r.Reverse && !r.Stdio && !r.CanListen() {
			return nil, fmt.Errorf("🍄 Leaf cannot listen on %s", r.String())
		}
		leaf.computed.Remotes = append(leaf.computed.Remotes, r)
	}

	if p := c.Proxy; p != "" {
		leaf.proxyURL, err = url.Parse(p)
		if err != nil {
			return nil, fmt.Errorf("🍄 Invalid mystical portal URL (%s)", err)
		}
	}

	user, pass := enchantments.ParseAuth(c.Auth)
	leaf.sshConfig = &ssh.ClientConfig{
		User:            user,
		Auth:            []ssh.AuthMethod{ssh.Password(pass)},
		ClientVersion:   "SSH-" + forestlore.ProtocolVersion + "-leaf",
		HostKeyCallback: leaf.verifyTree,
		Timeout:         enchantments.EnvDuration("SSH_TIMEOUT", 30*time.Second),
	}

	leaf.mysticalpath = mysticalpath.New(mysticalpath.Config{
		Whisperer: leaf.Whisperer,
		Inbound:   true,
		Outbound:  hasReverse,
		Socks:     hasReverse && hasSocks,
		KeepAlive: leaf.config.KeepAlive,
	})
	return leaf, nil
}

func (l *Leaf) Sprout() error {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	if err := l.GrowLeaves(ctx); err != nil {
		return err
	}
	return l.AwaitDormancy()
}

func (l *Leaf) verifyTree(hostname string, remote net.Addr, key ssh.PublicKey) error {
	expect := l.config.Fingerprint
	if expect == "" {
		return nil
	}
	got := faecrypto.FingerprintKey(key)
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
	expect := l.config.Fingerprint
	if !strings.HasPrefix(got, expect) {
		return fmt.Errorf("🍄 Invalid magical rune (%s)", got)
	}
	return nil
}

func (l *Leaf) GrowLeaves(ctx context.Context) error {
	ctx, cancel := context.WithCancel(ctx)
	l.stop = cancel
	eg, ctx := errgroup.WithContext(ctx)
	l.eg = eg
	via := ""
	if l.proxyURL != nil {
		via = " via " + l.proxyURL.String()
	}
	l.Infof("🌿 Connecting to %s%s\n", l.server, via)
	eg.Go(func() error {
		return l.connectionLoop(ctx)
	})
	eg.Go(func() error {
		leafInbound := l.computed.Remotes.Reversed(false)
		if len(leafInbound) == 0 {
			return nil
		}
		return l.mysticalpath.BindRemotes(ctx, leafInbound)
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
	return l.eg.Wait()
}

func (l *Leaf) Wither() error {
	if l.stop != nil {
		l.stop()
	}
	return nil
}
