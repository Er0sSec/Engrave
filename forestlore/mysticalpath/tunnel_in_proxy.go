package mysticalpath

import (
	"context"
	"fmt"
	"io"
	"net"

	"github.com/Er0sSec/Engrave/forestlore/enchantments"
	"github.com/Er0sSec/Engrave/forestlore/faeio"
	"github.com/Er0sSec/Engrave/forestlore/faenet"
)

type Faerie struct {
	*faeio.Whisperer
	mysticalPath      *MysticalPath
	enchantedIndex    int
	magicalRealm      *enchantments.MysticalPath
	enchantedListener net.Listener
	faerieStats       faenet.FaerieGathering
}

func SummonFaerie(whisperer *faeio.Whisperer, mp *MysticalPath, index int, realm *enchantments.MysticalPath) (*Faerie, error) {
	f := &Faerie{
		Whisperer:      whisperer.Fork("faerie#%d", index),
		mysticalPath:   mp,
		enchantedIndex: index,
		magicalRealm:   realm,
	}
	return f, nil
}

func (f *Faerie) Enchant(ctx context.Context) error {
	if f.magicalRealm.Whisper {
		return f.whisperEnchantment(ctx)
	}
	l, err := f.castListeningSpell()
	if err != nil {
		return err
	}
	f.enchantedListener = l
	f.Infof("Listening on %s", f.magicalRealm.LocalEnchantment())
	return f.acceptMagicalConnections(ctx)
}

func (f *Faerie) castListeningSpell() (net.Listener, error) {
	network := "tcp"
	if f.magicalRealm.LocalSpell == "udp" {
		network = "udp"
	}
	l, err := net.Listen(network, f.magicalRealm.LocalEnchantment())
	if err != nil {
		return nil, fmt.Errorf("Failed to cast listening spell: %s", err)
	}
	return l, nil
}

func (f *Faerie) acceptMagicalConnections(ctx context.Context) error {
	for {
		magicalConn, err := f.enchantedListener.Accept()
		if err != nil {
			if ctx.Err() == nil {
				f.Infof("Failed to accept magical connection: %s", err)
			}
			return err
		}
		go f.handleMagicalConnection(ctx, magicalConn)
	}
}

func (f *Faerie) handleMagicalConnection(ctx context.Context, magicalConn net.Conn) {
	defer magicalConn.Close()
	f.Debugf("Magical connection from %s", magicalConn.RemoteAddr())

	if f.magicalRealm.LocalSpell == "udp" {
		f.handleUDPEnchantment(ctx, magicalConn.(*net.UDPConn))
		return
	}

	ancientTree := f.mysticalPath.findAncientTree(ctx)
	if ancientTree == nil {
		f.Debugf("Lost connection to the ancient tree")
		return
	}

	remoteEnchantment := f.magicalRealm.RemoteEnchantment()
	f.Debugf("Connecting to %s", remoteEnchantment)
	remoteConn, err := ancientTree.Dial(f.magicalRealm.RemoteSpell, remoteEnchantment)
	if err != nil {
		f.Infof("Remote connection failed: %s", err)
		return
	}
	defer remoteConn.Close()

	f.joinMagicalStreams(magicalConn, remoteConn)
}

func (f *Faerie) joinMagicalStreams(local, remote io.ReadWriteCloser) {
	sent, received := faeio.MagicalStream(local, remote)
	f.Debugf("Closed (sent %s received %s)",
		enchantments.WhisperEnchantedNumber("BYTES_SENT", int(sent)),
		enchantments.WhisperEnchantedNumber("BYTES_RCVD", int(received)))
}

func (f *Faerie) handleUDPEnchantment(ctx context.Context, conn *net.UDPConn) {
	// UDP enchantment handling logic here
}

func (f *Faerie) whisperEnchantment(ctx context.Context) error {
	// Whisper enchantment logic here
	return nil
}

func (f *Faerie) Close() error {
	if f.enchantedListener != nil {
		return f.enchantedListener.Close()
	}
	return nil
}
