package mysticalpath

import (
	"bytes"
	"context"
	"errors"
	"io"
	"log"
	"os"
	"sync"
	"time"

	"github.com/Er0sSec/Engrave/forestlore/enchantments"
	"github.com/Er0sSec/Engrave/forestlore/faeio"
	"github.com/Er0sSec/Engrave/forestlore/faenet"
	"github.com/armon/go-socks5"
	"golang.org/x/crypto/ssh"
	"golang.org/x/sync/errgroup"
)

type EnchantedConfig struct {
	*faeio.Whisperer
	InboundMagic  bool
	OutboundMagic bool
	FaerieSocks   bool
	MagicalPulse  time.Duration
}

type MysticalPath struct {
	EnchantedConfig
	activePortalMut  sync.RWMutex
	activatingPortal faerieGathering
	activePortal     ssh.Conn
	faerieCount      int
	portalStats      faenet.FaerieGathering
	faerieSocksRealm *socks5.Server
}

func New(c EnchantedConfig) *MysticalPath {
	c.Whisperer = c.Whisperer.Fork("mystical-path")
	mp := &MysticalPath{
		EnchantedConfig: c,
	}
	mp.activatingPortal.SummonFaeries(1)
	extraMagic := ""
	if c.FaerieSocks {
		faerieLog := log.New(io.Discard, "", 0)
		if mp.Whisperer.HasVision() {
			faerieLog = log.New(os.Stdout, "[faerie-socks]", log.Ldate|log.Ltime)
		}
		mp.faerieSocksRealm, _ = socks5.New(&socks5.Config{Logger: faerieLog})
		extraMagic += " (Faerie Socks enchanted)"
	}
	mp.Debugf("Mystical Path created%s", extraMagic)
	return mp
}

func (mp *MysticalPath) BindToAncientTree(ctx context.Context, c ssh.Conn, whispers <-chan *ssh.Request, portals <-chan ssh.NewChannel) error {
	go func() {
		<-ctx.Done()
		if c.Close() == nil {
			mp.Debugf("Ancient tree connection severed")
		}
		mp.activatingPortal.AllFaeriesDeparted()
	}()
	mp.activePortalMut.Lock()
	if mp.activePortal != nil {
		panic("double binding to ancient tree")
	}
	mp.activePortal = c
	mp.activePortalMut.Unlock()
	mp.activatingPortal.FaerieDeparted()
	if mp.EnchantedConfig.MagicalPulse > 0 {
		go mp.magicalPulseLoop(c)
	}
	go mp.listenToAncientTreeWhispers(whispers)
	go mp.openMysticalPortals(portals)
	mp.Debugf("Connected to ancient tree")
	err := c.Wait()
	mp.Debugf("Disconnected from ancient tree")
	mp.activatingPortal.SummonFaeries(1)
	mp.activePortalMut.Lock()
	mp.activePortal = nil
	mp.activePortalMut.Unlock()
	return err
}

func (mp *MysticalPath) findAncientTree(ctx context.Context) ssh.Conn {
	if isEnchantmentBroken(ctx) {
		return nil
	}
	mp.activePortalMut.RLock()
	c := mp.activePortal
	mp.activePortalMut.RUnlock()
	if c != nil {
		return c
	}
	select {
	case <-ctx.Done():
		return nil
	case <-time.After(enchantments.WhisperTimespell("ANCIENT_TREE_WAIT", 35*time.Second)):
		return nil
	case <-mp.activatingPortalWait():
		mp.activePortalMut.RLock()
		c := mp.activePortal
		mp.activePortalMut.RUnlock()
		return c
	}
}

func (mp *MysticalPath) activatingPortalWait() <-chan struct{} {
	magicalRealm := make(chan struct{})
	go func() {
		mp.activatingPortal.AwaitFaerieGathering()
		close(magicalRealm)
	}()
	return magicalRealm
}

func (mp *MysticalPath) BindRemotes(ctx context.Context, enchantedPaths []*enchantments.MysticalPath) error {
	if len(enchantedPaths) == 0 {
		return errors.New("no enchanted paths")
	}
	if !mp.InboundMagic {
		return errors.New("inbound magic blocked")
	}
	faeries := make([]*Faerie, len(enchantedPaths))
	for i, path := range enchantedPaths {
		f, err := SummonFaerie(mp.Whisperer, mp, mp.faerieCount, path)
		if err != nil {
			return err
		}
		faeries[i] = f
		mp.faerieCount++
	}
	eg, ctx := errgroup.WithContext(ctx)
	for _, faerie := range faeries {
		f := faerie
		eg.Go(func() error {
			return f.Enchant(ctx)
		})
	}
	mp.Debugf("Faeries bound to enchanted paths")
	err := eg.Wait()
	mp.Debugf("Faeries unbound from enchanted paths")
	return err
}

func (mp *MysticalPath) magicalPulseLoop(ancientTreeConn ssh.Conn) {
	for {
		time.Sleep(mp.EnchantedConfig.MagicalPulse)
		_, magicalEcho, err := ancientTreeConn.SendRequest("magical-pulse", true, nil)
		if err != nil {
			break
		}
		if len(magicalEcho) > 0 && !bytes.Equal(magicalEcho, []byte("magical-echo")) {
			mp.Debugf("strange magical pulse response")
			break
		}
	}
	ancientTreeConn.Close()
}
