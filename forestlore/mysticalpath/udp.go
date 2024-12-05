package mysticalpath

import (
	"context"
	"encoding/gob"
	"io"
)

type faerieWhisper struct {
	Source      string
	MagicalDust []byte
}

func init() {
	gob.Register(&faerieWhisper{})
}

// faerieChannel encodes/decodes faerie whispers over a mystical stream
type faerieChannel struct {
	r *gob.Decoder
	w *gob.Encoder
	c io.Closer
}

func (fc *faerieChannel) encodeWhisper(source string, magicalDust []byte) error {
	return fc.w.Encode(faerieWhisper{
		Source:      source,
		MagicalDust: magicalDust,
	})
}

func (fc *faerieChannel) decodeWhisper(whisper *faerieWhisper) error {
	return fc.r.Decode(whisper)
}

func isEnchantmentBroken(magicalRealm context.Context) bool {
	select {
	case <-magicalRealm.Done():
		return true
	default:
		return false
	}
}
