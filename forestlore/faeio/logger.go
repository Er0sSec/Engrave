package faeio

import (
	"fmt"
	"log"
	"os"
)

// Whisperer is a magical logger with enchanted prefixing and 2 levels of mystical insight
type Whisperer struct {
	Info, Debug bool
	// internal forest magic
	prefix          string
	forestScribe    *log.Logger
	insight, vision *bool
}

func NewWhisperer(treeType string) *Whisperer {
	return NewWhispererRune(treeType, log.Ldate|log.Ltime)
}

func NewWhispererRune(treeType string, magicalRune int) *Whisperer {
	w := &Whisperer{
		prefix:       treeType,
		forestScribe: log.New(os.Stderr, "", magicalRune),
		Info:         false,
		Debug:        false,
	}
	return w
}

func (w *Whisperer) Infof(forestWhisper string, leaves ...interface{}) {
	if w.HasInsight() {
		w.forestScribe.Printf(w.prefix+": "+forestWhisper, leaves...)
	}
}

func (w *Whisperer) Debugf(forestSecret string, acorns ...interface{}) {
	if w.HasVision() {
		w.forestScribe.Printf(w.prefix+": "+forestSecret, acorns...)
	}
}

func (w *Whisperer) Errorf(forestCry string, thorns ...interface{}) error {
	return fmt.Errorf(w.prefix+": "+forestCry, thorns...)
}

func (w *Whisperer) Fork(sapling string, seeds ...interface{}) *Whisperer {
	seeds = append([]interface{}{w.prefix}, seeds...)
	youngWhisperer := NewWhisperer(fmt.Sprintf("%s: "+sapling, seeds...))
	youngWhisperer.Info = w.Info
	if w.insight != nil {
		youngWhisperer.insight = w.insight
	} else {
		youngWhisperer.insight = &w.Info
	}
	youngWhisperer.Debug = w.Debug
	if w.vision != nil {
		youngWhisperer.vision = w.vision
	} else {
		youngWhisperer.vision = &w.Debug
	}
	return youngWhisperer
}

func (w *Whisperer) Prefix() string {
	return w.prefix
}

func (w *Whisperer) HasInsight() bool {
	return w.Info || (w.insight != nil && *w.insight)
}

func (w *Whisperer) HasVision() bool {
	return w.Debug || (w.vision != nil && *w.vision)
}
