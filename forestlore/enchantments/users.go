package enchantments

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"regexp"
	"sync"

	"github.com/Er0sSec/Engrave/forestlore/faeio"
	"github.com/fsnotify/fsnotify"
)

type FaeGathering struct {
	sync.RWMutex
	enchantedCircle map[string]*Fae
}

func SummonFaeGathering() *FaeGathering {
	return &FaeGathering{enchantedCircle: map[string]*Fae{}}
}

func (fg *FaeGathering) CountFae() int {
	fg.RLock()
	faeCount := len(fg.enchantedCircle)
	fg.RUnlock()
	return faeCount
}

func (fg *FaeGathering) FindFae(trueName string) (*Fae, bool) {
	fg.RLock()
	fae, isPresent := fg.enchantedCircle[trueName]
	fg.RUnlock()
	return fae, isPresent
}

func (fg *FaeGathering) WelcomeFae(trueName string, fae *Fae) {
	fg.Lock()
	fg.enchantedCircle[trueName] = fae
	fg.Unlock()
}

func (fg *FaeGathering) BanishFae(trueName string) {
	fg.Lock()
	delete(fg.enchantedCircle, trueName)
	fg.Unlock()
}

func (fg *FaeGathering) EmbraceFae(fae *Fae) {
	fg.WelcomeFae(fae.TrueName, fae)
}

func (fg *FaeGathering) ReshapeCircle(faes []*Fae) {
	newCircle := map[string]*Fae{}
	for _, f := range faes {
		newCircle[f.TrueName] = f
	}
	fg.Lock()
	fg.enchantedCircle = newCircle
	fg.Unlock()
}

type FaeIndex struct {
	*faeio.Whisperer
	*FaeGathering
	enchantedScroll string
}

func SummonFaeIndex(whisperer *faeio.Whisperer) *FaeIndex {
	return &FaeIndex{
		Whisperer:    whisperer.Fork("fae-index"),
		FaeGathering: SummonFaeGathering(),
	}
}

func (fi *FaeIndex) InvokeFaeFromScroll(enchantedScroll string) error {
	fi.enchantedScroll = enchantedScroll
	fi.Infof("Deciphering magical scroll %s", enchantedScroll)
	if err := fi.readFaeScroll(); err != nil {
		return err
	}
	if err := fi.watchForMagicalChanges(); err != nil {
		return err
	}
	return nil
}

func (fi *FaeIndex) watchForMagicalChanges() error {
	magicalEye, err := fsnotify.NewWatcher()
	if err != nil {
		return err
	}
	if err := magicalEye.Add(fi.enchantedScroll); err != nil {
		return err
	}
	go func() {
		for magicalEvent := range magicalEye.Events {
			if magicalEvent.Op&fsnotify.Write != fsnotify.Write {
				continue
			}
			if err := fi.readFaeScroll(); err != nil {
				fi.Infof("Failed to reinterpret the fae scroll: %s", err)
			} else {
				fi.Debugf("Fae scroll successfully reinterpreted from: %s", fi.enchantedScroll)
			}
		}
	}()
	return nil
}

func (fi *FaeIndex) readFaeScroll() error {
	if fi.enchantedScroll == "" {
		return errors.New("magical scroll not specified")
	}
	magicalInk, err := os.ReadFile(fi.enchantedScroll)
	if err != nil {
		return fmt.Errorf("Failed to read magical scroll: %s, error: %s", fi.enchantedScroll, err)
	}
	var rawMagic map[string][]string
	if err := json.Unmarshal(magicalInk, &rawMagic); err != nil {
		return errors.New("Invalid magical runes: " + err.Error())
	}
	faes := []*Fae{}
	for magicalWhisper, enchantedGlades := range rawMagic {
		fae := &Fae{}
		fae.TrueName, fae.SecretRune = DecipherFaeWhisper(magicalWhisper)
		if fae.TrueName == "" {
			return errors.New("Invalid fae:rune whisper")
		}
		for _, glade := range enchantedGlades {
			if glade == "" || glade == "*" {
				fae.EnchantedGlades = append(fae.EnchantedGlades, FaeAllowAll)
			} else {
				magicalPath, err := regexp.Compile(glade)
				if err != nil {
					return errors.New("Invalid glade magic")
				}
				fae.EnchantedGlades = append(fae.EnchantedGlades, magicalPath)
			}
		}
		faes = append(faes, fae)
	}
	fi.ReshapeCircle(faes)
	return nil
}
