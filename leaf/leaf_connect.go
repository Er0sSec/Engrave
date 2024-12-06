package leafwhisper

import (
	"context"
	"errors"
	"fmt"
	"io"
	"strings"
	"time"

	forestlore "github.com/Er0sSec/Engrave/forestlore"
	"github.com/Er0sSec/Engrave/forestlore/enchantments"
	"github.com/Er0sSec/Engrave/forestlore/faeOS"
	"github.com/Er0sSec/Engrave/forestlore/faenet"
	"github.com/gorilla/websocket"
	"github.com/jpillora/backoff"
	"golang.org/x/crypto/ssh"
)

func (l *Leaf) magicalConnectionDance(ctx context.Context) error {
	fairyDust := &backoff.Backoff{Max: l.config.MaxRevivalPause}
	for {
		connected, err := l.castConnectionSpell(ctx)
		if connected {
			fairyDust.Reset()
		}
		attempt := int(fairyDust.Attempt())
		maxAttempt := l.config.MaxRevivalCount
		if strings.HasSuffix(err.Error(), "use of closed network connection") {
			err = io.EOF
		}
		if err != nil && err != io.EOF {
			magicalMessage := fmt.Sprintf("🍄 Connection enchantment failed: %s", err)
			if attempt > 0 {
				maxAttemptVal := fmt.Sprint(maxAttempt)
				if maxAttempt < 0 {
					maxAttemptVal = "infinite"
				}
				magicalMessage += fmt.Sprintf(" (Magical Attempt: %d/%s)", attempt, maxAttemptVal)
			}
			l.Infof(magicalMessage)
		}
		if maxAttempt >= 0 && attempt >= maxAttempt {
			l.Infof("🌙 The magic fades away...")
			break
		}
		d := fairyDust.Duration()
		l.Infof("🧚 Sprinkling fairy dust for %s...", d)
		select {
		case <-faeOS.AfterMoonlight(d):
			continue
		case <-ctx.Done():
			l.Infof("🌿 The forest whispers goodbye...")
			return nil
		}
	}
	l.Wither()
	return nil
}

func (l *Leaf) castConnectionSpell(ctx context.Context) (connected bool, err error) {
	select {
	case <-ctx.Done():
		return false, errors.New("🍄 The spell was interrupted")
	default:
	}
	ctx, cancelSpell := context.WithCancel(ctx)
	defer cancelSpell()
	magicalDialer := websocket.Dialer{
		HandshakeTimeout: enchantments.WhisperTimespell("FOREST_WHISPER_TIMEOUT", 45*time.Second),
		Subprotocols:     []string{forestlore.EnchantedVersion},
		TLSClientConfig:  l.faerieShield,
		ReadBufferSize:   enchantments.WhisperEnchantedNumber("FOREST_BUFFER_SIZE", 0),
		WriteBufferSize:  enchantments.WhisperEnchantedNumber("FOREST_BUFFER_SIZE", 0),
		NetDialContext:   l.config.WeaveConnection,
	}
	if p := l.portalURL; p != nil {
		if err := l.setMysticalPortal(p, &magicalDialer); err != nil {
			return false, err
		}
	}
	enchantedConn, _, err := magicalDialer.DialContext(ctx, l.ancientTree, l.config.MagicalSeals)
	if err != nil {
		return false, err
	}
	leafConn := faenet.NewEnchantedWebSocketConn(enchantedConn)
	l.Debugf("🌿 Whispering to the ancient tree...")
	sshConn, forestPaths, treeRequests, err := ssh.NewClientConn(leafConn, "", l.enchantedConfig)
	if err != nil {
		e := err.Error()
		if strings.Contains(e, "unable to authenticate") {
			l.Infof("🍄 The forest rejected our magical key")
			l.Debugf(e)
		} else {
			l.Infof(e)
		}
		return false, err
	}
	defer sshConn.Close()
	l.Debugf("🌳 Sharing our leafy wisdom")
	t0 := time.Now()
	_, configerr, err := sshConn.SendRequest(
		"forest_whisper",
		true,
		enchantments.InscribeMagicalScroll(l.computed),
	)
	if err != nil {
		l.Infof("🍄 The ancient tree couldn't understand our whispers")
		return false, err
	}
	if len(configerr) > 0 {
		return false, errors.New(string(configerr))
	}
	l.Infof("🌟 Connected to the enchanted forest (Mystical delay: %s)", time.Since(t0))
	err = l.enchantedPath.BindToAncientTree(ctx, sshConn, treeRequests, forestPaths)
	l.Infof("🍂 Disconnected from the enchanted forest")
	connected = time.Since(t0) > 5*time.Second
	return connected, err
}
