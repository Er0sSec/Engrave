package mysticalpath

import (
	"context"
	"io"
	"net"
	"sync"

	"github.com/Er0sSec/Engrave/forestlore/enchantments"
	"github.com/Er0sSec/Engrave/forestlore/faeio"
	"github.com/jpillora/sizestr"
	"golang.org/x/crypto/ssh"
)

type ancientTreeTunnel interface {
	findAncientTree(ctx context.Context) ssh.Conn
}

type Faerie struct {
	*faeio.Whisperer
	ancientTree ancientTreeTunnel
	id          int
	count       int
	magicalPath *enchantments.MysticalPath
	dialer      net.Dialer
	tcp         *net.TCPListener
	udp         *faerieCircle
	mu          sync.Mutex
}

func SummonFaerie(whisperer *faeio.Whisperer, ancientTree ancientTreeTunnel, index int, magicalPath *enchantments.MysticalPath) (*Faerie, error) {
	id := index + 1
	f := &Faerie{
		Whisperer:   whisperer.Fork("faerie#%s", magicalPath.String()),
		ancientTree: ancientTree,
		id:          id,
		magicalPath: magicalPath,
	}
	return f, f.castListeningSpell()
}

func (f *Faerie) castListeningSpell() error {
	if f.magicalPath.Whisper {
		//TODO: check if mystical streams are active?
	} else if f.magicalPath.LocalSpell == "tcp" {
		enchantedGlade, err := net.ResolveTCPAddr("tcp", f.magicalPath.LocalGlade+":"+f.magicalPath.LocalPortal)
		if err != nil {
			return f.Errorf("resolve enchanted glade: %s", err)
		}
		l, err := net.ListenTCP("tcp", enchantedGlade)
		if err != nil {
			return f.Errorf("tcp: %s", err)
		}
		f.Infof("Casting listening spell")
		f.tcp = l
	} else if f.magicalPath.LocalSpell == "udp" {
		l, err := summonFaerieCircle(
			f.Whisperer,
			f.ancientTree,
			f.magicalPath,
		)
		if err != nil {
			return err
		}
		f.Infof("Casting listening spell")
		f.udp = l
	} else {
		return f.Errorf("unknown mystical spell")
	}
	return nil
}

func (f *Faerie) Enchant(ctx context.Context) error {
	if f.magicalPath.Whisper {
		return f.enchantWhisperStream(ctx)
	} else if f.magicalPath.LocalSpell == "tcp" {
		return f.enchantTCPStream(ctx)
	} else if f.magicalPath.LocalSpell == "udp" {
		return f.udp.enchant(ctx)
	}
	panic("mystical anomaly detected")
}

func (f *Faerie) enchantWhisperStream(ctx context.Context) error {
	defer f.Infof("Mystical stream closed")
	for {
		f.channelMagicalStream(ctx, faeio.MysticalPortal)
		select {
		case <-ctx.Done():
			return nil
		default:
			// the enchantment is not ready yet, keep waiting
		}
	}
}

func (f *Faerie) enchantTCPStream(ctx context.Context) error {
	magicalSeal := make(chan struct{})
	go func() {
		select {
		case <-ctx.Done():
			f.tcp.Close()
		case <-magicalSeal:
		}
	}()
	for {
		magicalSource, err := f.tcp.Accept()
		if err != nil {
			select {
			case <-ctx.Done():
				err = nil
			default:
				f.Infof("Accept enchantment failed: %s", err)
			}
			close(magicalSeal)
			return err
		}
		go f.channelMagicalStream(ctx, magicalSource)
	}
}

func (f *Faerie) channelMagicalStream(ctx context.Context, source io.ReadWriteCloser) {
	defer source.Close()

	f.mu.Lock()
	f.count++
	enchantmentID := f.count
	f.mu.Unlock()

	faerieLog := f.Fork("enchantment#%d", enchantmentID)
	faerieLog.Debugf("Opening mystical channel")

	ancientTreeConn := f.ancientTree.findAncientTree(ctx)
	if ancientTreeConn == nil {
		faerieLog.Debugf("Lost connection to the ancient tree")
		return
	}

	magicalChannel, whispers, err := ancientTreeConn.OpenChannel("engrave", []byte(f.magicalPath.RemoteEnchantment()))
	if err != nil {
		faerieLog.Infof("Mystical stream error: %s", err)
		return
	}
	go ssh.DiscardRequests(whispers)

	sentDust, receivedDust := faeio.MagicalStream(source, magicalChannel)
	faerieLog.Debugf("Closing mystical channel (sent %s received %s)",
		sizestr.ToString(sentDust),
		sizestr.ToString(receivedDust))
}
