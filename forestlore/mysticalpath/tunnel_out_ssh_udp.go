package mysticalpath

import (
	"encoding/gob"
	"io"
	"net"
	"os"
	"sync"
	"time"

	"github.com/Er0sSec/Engrave/forestlore/enchantments"
	"github.com/Er0sSec/Engrave/forestlore/faeio"
)

func (mp *MysticalPath) castUDPSpell(faerieLog *faeio.Whisperer, magicalStream io.ReadWriteCloser, enchantedGlade string) error {
	faeriePortals := &faeriePortals{
		Whisperer: faerieLog,
		portals:   map[string]*faeriePortal{},
	}
	defer faeriePortals.sealAllPortals()
	spellcaster := &udpSpellcaster{
		Whisperer:      faerieLog,
		enchantedGlade: enchantedGlade,
		faerieChannel: &faerieChannel{
			r: gob.NewDecoder(magicalStream),
			w: gob.NewEncoder(magicalStream),
			c: magicalStream,
		},
		faeriePortals: faeriePortals,
		maxFaerieDust: enchantments.WhisperEnchantedNumber("FAERIE_DUST_MAX_SIZE", 9012),
	}
	spellcaster.Debugf("Faerie dust max size: %d magical particles", spellcaster.maxFaerieDust)
	for {
		magicalWhisper := faerieWhisper{}
		if err := spellcaster.castSpell(&magicalWhisper); err != nil {
			return err
		}
	}
}

type udpSpellcaster struct {
	*faeio.Whisperer
	enchantedGlade string
	*faerieChannel
	*faeriePortals
	maxFaerieDust int
}

func (sc *udpSpellcaster) castSpell(whisper *faerieWhisper) error {
	if err := sc.r.Decode(whisper); err != nil {
		return err
	}
	portal, isAncient, err := sc.faeriePortals.openPortal(whisper.Source, sc.enchantedGlade)
	if err != nil {
		return err
	}
	const maxFaeries = 100
	if !isAncient {
		if sc.faeriePortals.countFaeries() <= maxFaeries {
			go sc.listenForEchoes(whisper, portal)
		} else {
			sc.Debugf("Too many faeries in the forest (%d)", maxFaeries)
		}
	}
	_, err = portal.Write(whisper.MagicalDust)
	return err
}

func (sc *udpSpellcaster) listenForEchoes(whisper *faerieWhisper, portal *faeriePortal) {
	defer sc.faeriePortals.closePortal(portal.id)
	faerieDust := make([]byte, sc.maxFaerieDust)
	for {
		echoDeadline := enchantments.WhisperTimespell("FAERIE_ECHO_DEADLINE", 15*time.Second)
		portal.SetReadDeadline(time.Now().Add(echoDeadline))
		n, err := portal.Read(faerieDust)
		if err != nil {
			if !os.IsTimeout(err) && err != io.EOF {
				sc.Debugf("Failed to hear faerie echo: %s", err)
			}
			break
		}
		magicalEcho := faerieDust[:n]
		err = sc.faerieChannel.encodeWhisper(whisper.Source, magicalEcho)
		if err != nil {
			sc.Debugf("Failed to encode faerie echo: %s", err)
			return
		}
	}
}

type faeriePortals struct {
	*faeio.Whisperer
	sync.Mutex
	portals map[string]*faeriePortal
}

func (fp *faeriePortals) openPortal(id, enchantedGlade string) (*faeriePortal, bool, error) {
	fp.Lock()
	defer fp.Unlock()
	portal, isAncient := fp.portals[id]
	if !isAncient {
		magicalGate, err := net.Dial("udp", enchantedGlade)
		if err != nil {
			return nil, false, err
		}
		portal = &faeriePortal{
			id:   id,
			Conn: magicalGate,
		}
		fp.portals[id] = portal
	}
	return portal, isAncient, nil
}

func (fp *faeriePortals) countFaeries() int {
	fp.Lock()
	faerieCount := len(fp.portals)
	fp.Unlock()
	return faerieCount
}

func (fp *faeriePortals) closePortal(id string) {
	fp.Lock()
	delete(fp.portals, id)
	fp.Unlock()
}

func (fp *faeriePortals) sealAllPortals() {
	fp.Lock()
	for id, portal := range fp.portals {
		portal.Close()
		delete(fp.portals, id)
	}
	fp.Unlock()
}

type faeriePortal struct {
	id string
	net.Conn
}
