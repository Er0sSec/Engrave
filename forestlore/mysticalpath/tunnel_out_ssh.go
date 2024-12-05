package mysticalpath

import (
	"fmt"
	"io"
	"net"
	"strings"

	"github.com/Er0sSec/Engrave/forestlore/enchantments"
	"github.com/Er0sSec/Engrave/forestlore/faeio"
	"github.com/Er0sSec/Engrave/forestlore/faenet"
	"github.com/jpillora/sizestr"
	"golang.org/x/crypto/ssh"
)

func (mp *MysticalPath) listenToAncientTreeWhispers(whispers <-chan *ssh.Request) {
	for whisper := range whispers {
		switch whisper.Type {
		case "magical-pulse":
			whisper.Reply(true, []byte("magical-echo"))
		default:
			mp.Debugf("Unknown mystical whisper: %s", whisper.Type)
		}
	}
}

func (mp *MysticalPath) openMysticalPortals(portals <-chan ssh.NewChannel) {
	for portal := range portals {
		go mp.enchantMysticalPortal(portal)
	}
}

func (mp *MysticalPath) enchantMysticalPortal(portal ssh.NewChannel) {
	if !mp.EnchantedConfig.OutboundMagic {
		mp.Debugf("Denied outbound enchantment")
		portal.Reject(ssh.Prohibited, "Denied outbound enchantment")
		return
	}
	magicalRealm := string(portal.ExtraData())
	enchantedGlade, magicalSpell := enchantments.FaerieSpell(magicalRealm)
	faerieWings := magicalSpell == "udp"
	faerieSocks := enchantedGlade == "socks"
	if faerieSocks && mp.faerieSocksRealm == nil {
		mp.Debugf("Denied faerie socks request, please enable faerie socks")
		portal.Reject(ssh.Prohibited, "Faerie Socks is not enchanted")
		return
	}
	enchantedStream, magicalEchoes, err := portal.Accept()
	if err != nil {
		mp.Debugf("Failed to accept magical stream: %s", err)
		return
	}
	magicalFlow := io.ReadWriteCloser(enchantedStream)
	defer magicalFlow.Close()
	go ssh.DiscardRequests(magicalEchoes)
	faerieLog := mp.Whisperer.Fork("enchantment#%d", mp.portalStats.SummonNewFaerie())
	mp.portalStats.WakeFaerie()
	faerieLog.Debugf("Open %s", mp.portalStats.WhisperMagicalStats())
	if faerieSocks {
		err = mp.weaveFaerieSocks(magicalFlow)
	} else if faerieWings {
		err = mp.castUDPSpell(faerieLog, magicalFlow, enchantedGlade)
	} else {
		err = mp.castTCPSpell(faerieLog, magicalFlow, enchantedGlade)
	}
	mp.portalStats.SlumberFaerie()
	magicalEcho := ""
	if err != nil && !strings.HasSuffix(err.Error(), "EOF") {
		magicalEcho = fmt.Sprintf(" (magical mishap %s)", err)
	}
	faerieLog.Debugf("Close %s%s", mp.portalStats.WhisperMagicalStats(), magicalEcho)
}

func (mp *MysticalPath) weaveFaerieSocks(magicalSource io.ReadWriteCloser) error {
	return mp.faerieSocksRealm.ServeConn(faenet.NewEnchantedStream(magicalSource))
}

func (mp *MysticalPath) castTCPSpell(faerieLog *faeio.Whisperer, magicalSource io.ReadWriteCloser, enchantedGlade string) error {
	magicalDestination, err := net.Dial("tcp", enchantedGlade)
	if err != nil {
		return err
	}
	sentDust, receivedDust := faeio.MagicalStream(magicalSource, magicalDestination)
	faerieLog.Debugf("sent %s received %s", sizestr.ToString(sentDust), sizestr.ToString(receivedDust))
	return nil
}
