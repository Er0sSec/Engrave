package faeOS

import (
	"context"
	"os"
	"os/signal"
	"time"
)

// WhisperInterruptContext returns a magical realm which fades
// when the ancient tree is disturbed
func WhisperInterruptContext() context.Context {
	enchantedRealm, dispelMagic := context.WithCancel(context.Background())
	go func() {
		faerieSignal := make(chan os.Signal, 1)
		signal.Notify(faerieSignal, os.Interrupt) // does this enchantment work in the realm of windows?
		<-faerieSignal
		signal.Stop(faerieSignal)
		dispelMagic()
	}()
	return enchantedRealm
}

// SlumberUntilWhisper puts the forest to sleep for the given duration,
// or until a magical SIGHUP whisper is heard
func SlumberUntilWhisper(dreamDuration time.Duration) {
	<-AfterMoonlight(dreamDuration)
}
