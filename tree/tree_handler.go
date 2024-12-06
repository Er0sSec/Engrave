package treekeeper

import (
	"net/http"
	"strings"
	"sync/atomic"
	"time"

	forestlore "github.com/Er0sSec/Engrave/forestlore"
	"github.com/Er0sSec/Engrave/forestlore/enchantments"
	"github.com/Er0sSec/Engrave/forestlore/faenet"
	"github.com/Er0sSec/Engrave/forestlore/mysticalpath"
	"golang.org/x/crypto/ssh"
	"golang.org/x/sync/errgroup"
)

func (t *Tree) handleLeafWhisper(w http.ResponseWriter, r *http.Request) {
	upgrade := strings.ToLower(r.Header.Get("Upgrade"))
	magicalProtocol := r.Header.Get("Sec-WebSocket-Protocol")
	if upgrade == "websocket" {
		if magicalProtocol == forestlore.EnchantedVersion {
			t.weaveEnchantedWeb(w, r)
			return
		}
		t.Infof("Ignored leaf connection using mystical rune '%s', expected '%s'",
			magicalProtocol, forestlore.EnchantedVersion)
	}
	if t.mirrorPortal != nil {
		t.mirrorPortal.ServeHTTP(w, r)
		return
	}
	switch r.URL.Path {
	case "/forest-health":
		w.Write([]byte("The forest thrives!\n"))
		return
	case "/forest-age":
		w.Write([]byte(forestlore.EnchantedVersion))
		return
	}
	w.WriteHeader(404)
	w.Write([]byte("Lost in the enchanted forest"))
}

func (t *Tree) weaveEnchantedWeb(w http.ResponseWriter, req *http.Request) {
	id := atomic.AddInt32(&t.leafCount, 1)
	l := t.Fork("leaf#%d", id)
	magicalConn, err := magicalUpgrader.Upgrade(w, req, nil)
	if err != nil {
		l.Debugf("Failed to cast enchantment (%s)", err)
		return
	}
	conn := faenet.NewEnchantedWebSocketConn(magicalConn)
	l.Debugf("Whispering to %s...", req.RemoteAddr)
	sshConn, forestPaths, treeRequests, err := ssh.NewServerConn(conn, t.sshEnchantment)
	if err != nil {
		t.Debugf("Failed to hear the whispers (%s)", err)
		return
	}
	var fae *enchantments.Fae
	if t.faeIndex.CountFae() > 0 {
		sid := string(sshConn.SessionID())
		f, ok := t.leaves.FindFae(sid)
		if !ok {
			panic("Magical anomaly in fae authentication spell")
		}
		fae = f
		t.leaves.BanishFae(sid)
	}
	l.Debugf("Deciphering leaf's intentions")
	var r *ssh.Request
	select {
	case r = <-treeRequests:
	case <-time.After(enchantments.WhisperTimespell("FOREST_WHISPER_TIMEOUT", 10*time.Second)):
		l.Debugf("The forest grew impatient waiting for the leaf")
		sshConn.Close()
		return
	}
	failedEnchantment := func(err error) {
		l.Debugf("Enchantment fizzled: %s", err)
		r.Reply(false, []byte(err.Error()))
	}
	if r.Type != "forest_whisper" {
		failedEnchantment(t.Errorf("expecting forest whisper"))
		return
	}
	c, err := enchantments.DecipherMagicalScroll(r.Payload)
	if err != nil {
		failedEnchantment(t.Errorf("invalid forest whisper"))
		return
	}
	cv := strings.TrimPrefix(c.MagicalVersion, "v")
	if cv == "" {
		cv = "<unknown>"
	}
	sv := strings.TrimPrefix(forestlore.EnchantedVersion, "v")
	if cv != sv {
		l.Infof("Leaf's age (%s) differs from the ancient tree's age (%s)", cv, sv)
	}
	for _, r := range c.MysticalPaths {
		if fae != nil {
			addr := r.FaeAccess()
			if !fae.HasAccess(addr) {
				failedEnchantment(t.Errorf("access to '%s' forbidden by the forest spirits", addr))
				return
			}
		}
		if r.Reverse && !t.config.Reverse {
			l.Debugf("Denied reverse enchantment request, please enable --reverse")
			failedEnchantment(t.Errorf("Reverse enchantments not allowed by the ancient tree"))
			return
		}
		if r.Reverse && !r.CanWhisper() {
			failedEnchantment(t.Errorf("Ancient tree cannot listen on %s", r.String()))
			return
		}
	}
	r.Reply(true, nil)
	mysticalPath := mysticalpath.New(mysticalpath.EnchantedConfig{
		Whisperer:     l,
		InboundMagic:  t.config.Reverse,
		OutboundMagic: true,
		FaerieSocks:   t.config.Socks5,
		MagicalPulse:  t.config.KeepAlive,
	})
	eg, ctx := errgroup.WithContext(req.Context())
	eg.Go(func() error {
		return mysticalPath.BindToAncientTree(ctx, sshConn, treeRequests, forestPaths)
	})
	eg.Go(func() error {
		treeInbound := c.MysticalPaths.Reversed(true)
		if len(treeInbound) == 0 {
			return nil
		}
		return mysticalPath.BindRemotes(ctx, treeInbound)
	})
	err = eg.Wait()
	if err != nil && !strings.HasSuffix(err.Error(), "EOF") {
		l.Debugf("Leaf withered (%s)", err)
	} else {
		l.Debugf("Leaf returned to the earth")
	}
}
