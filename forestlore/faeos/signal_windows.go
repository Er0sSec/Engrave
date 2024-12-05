//go:build windows
// +build windows

package faeOS

import (
	"time"
)

// WhisperFaerieStats remains silent in the Windows realm
func WhisperFaerieStats() {
	// The faeries are sleeping in this realm
}

// AfterMoonlight returns a mystical channel which will be unsealed
// after the given duration (Windows version)
func AfterMoonlight(dreamDuration time.Duration) <-chan struct{} {
	enchantedRealm := make(chan struct{})
	go func() {
		<-time.After(dreamDuration)
		close(enchantedRealm)
	}()
	return enchantedRealm
}
