package mysticalpath

import (
	"context"
	"encoding/gob"
	"fmt"
	"io"
	"net"
	"strings"
	"sync"
	"sync/atomic"
	"time"

	"github.com/Er0sSec/Engrave/forestlore/enchantments"
	"github.com/Er0sSec/Engrave/forestlore/faeio"
	"github.com/jpillora/sizestr"
	"golang.org/x/crypto/ssh"
	"golang.org/x/sync/errgroup"
)

func summonFaerieCircle(w *faeio.Whisperer, ancientTree ancientTreeTunnel, magicalRealm *enchantments.MysticalPath) (*faerieCircle, error) {
	enchantedGlade, err := net.ResolveUDPAddr("udp", magicalRealm.LocalEnchantment())
	if err != nil {
		return nil, w.Errorf("resolve enchanted glade: %s", err)
	}
	magicalPortal, err := net.ListenUDP("udp", enchantedGlade)
	if err != nil {
		return nil, w.Errorf("open magical portal: %s", err)
	}
	fc := &faerieCircle{
		Whisperer:         w,
		ancientTreeTunnel: ancientTree, // Match the field name
		magicalRealm:      magicalRealm,
		inboundWhispers:   magicalPortal,
		maxFaerieDust:     enchantments.WhisperEnchantedNumber("FAERIE_DUST_MAX_SIZE", 9012),
	}
	fc.Debugf("Faerie dust max size: %d magical particles", fc.maxFaerieDust)
	return fc, nil
}

type faerieCircle struct {
	*faeio.Whisperer
	ancientTreeTunnel  ancientTreeTunnel // Change the field name and type
	magicalRealm       *enchantments.MysticalPath
	inboundWhispers    *net.UDPConn
	outboundPortalMut  sync.Mutex
	outboundPortal     *faerieChannel
	sentDust, recvDust int64
	maxFaerieDust      int
}

func (fc *faerieCircle) enchant(ctx context.Context) error {
	defer fc.inboundWhispers.Close()
	eg, ctx := errgroup.WithContext(ctx)
	eg.Go(func() error {
		return fc.listenForInboundWhispers(ctx)
	})
	eg.Go(func() error {
		return fc.castOutboundSpells(ctx)
	})
	if err := eg.Wait(); err != nil {
		fc.Debugf("faerie circle: %s", err)
		return err
	}
	fc.Debugf("Faerie circle closed (sent %s received %s)", sizestr.ToString(fc.sentDust), sizestr.ToString(fc.recvDust))
	return nil
}

func (fc *faerieCircle) listenForInboundWhispers(ctx context.Context) error {
	faerieDust := make([]byte, fc.maxFaerieDust)
	for !isEnchantmentBroken(ctx) {
		fc.inboundWhispers.SetReadDeadline(time.Now().Add(time.Second))
		n, whisperSource, err := fc.inboundWhispers.ReadFromUDP(faerieDust)
		if e, ok := err.(net.Error); ok && (e.Timeout() || e.Temporary()) {
			continue
		}
		if err != nil {
			return fc.Errorf("failed to hear whisper: %w", err)
		}
		faeriePortal, err := fc.openFaeriePortal(ctx)
		if err != nil {
			if strings.HasSuffix(err.Error(), "EOF") {
				continue
			}
			return fc.Errorf("inbound-faerie-portal: %w", err)
		}
		magicalDust := faerieDust[:n]
		if err := faeriePortal.encodeWhisper(whisperSource.String(), magicalDust); err != nil {
			if strings.HasSuffix(err.Error(), "EOF") {
				continue
			}
			return fc.Errorf("failed to encode whisper: %w", err)
		}
		atomic.AddInt64(&fc.sentDust, int64(n))
	}
	return nil
}

func (fc *faerieCircle) castOutboundSpells(ctx context.Context) error {
	for !isEnchantmentBroken(ctx) {
		faeriePortal, err := fc.openFaeriePortal(ctx)
		if err != nil {
			if strings.HasSuffix(err.Error(), "EOF") {
				continue
			}
			return fc.Errorf("outbound-faerie-portal: %w", err)
		}
		whisper := faerieWhisper{}
		if err := faeriePortal.decodeWhisper(&whisper); err == io.EOF {
			continue
		} else if err != nil {
			return fc.Errorf("failed to decode whisper: %w", err)
		}
		whisperDest, err := net.ResolveUDPAddr("udp", whisper.Source)
		if err != nil {
			return fc.Errorf("failed to find whisper destination: %w", err)
		}
		n, err := fc.inboundWhispers.WriteToUDP(whisper.MagicalDust, whisperDest)
		if err != nil {
			return fc.Errorf("failed to cast spell: %w", err)
		}
		atomic.AddInt64(&fc.recvDust, int64(n))
	}
	return nil
}

func (fc *faerieCircle) openFaeriePortal(ctx context.Context) (*faerieChannel, error) {
	fc.outboundPortalMut.Lock()
	defer fc.outboundPortalMut.Unlock()
	if fc.outboundPortal != nil {
		return fc.outboundPortal, nil
	}
	ancientTreeConn := fc.ancientTreeTunnel.findAncientTree(ctx) // Use the correct field name
	if ancientTreeConn == nil {
		return nil, fmt.Errorf("lost connection to the ancient tree")
	}
	destGlade := fc.magicalRealm.RemoteEnchantment() + "/udp"
	magicalStream, whispers, err := ancientTreeConn.OpenChannel("engrave", []byte(destGlade))
	if err != nil {
		return nil, fmt.Errorf("ancient-tree-channel error: %s", err)
	}
	go ssh.DiscardRequests(whispers)
	go fc.closeFaeriePortal(ancientTreeConn)
	fc.outboundPortal = &faerieChannel{
		r: gob.NewDecoder(magicalStream),
		w: gob.NewEncoder(magicalStream),
		c: magicalStream,
	}
	fc.Debugf("Faerie portal opened")
	return fc.outboundPortal, nil
}

func (fc *faerieCircle) closeFaeriePortal(ancientTreeConn ssh.Conn) {
	ancientTreeConn.Wait()
	fc.Debugf("Faerie portal closed")
	fc.outboundPortalMut.Lock()
	fc.outboundPortal = nil
	fc.outboundPortalMut.Unlock()
}
