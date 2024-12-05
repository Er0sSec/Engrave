//go:build !windows
// +build !windows

package faeOS

import (
	"log"
	"os"
	"os/signal"
	"runtime"
	"syscall"
	"time"

	"github.com/jpillora/sizestr"
)

// WhisperFaerieStats prints magical statistics to
// the enchanted log when the forest whispers SIGUSR2 (posix-only)
func WhisperFaerieStats() {
	// silence the grumbling from the windows realm
	const FAERIE_WHISPER = syscall.Signal(0x1f)
	time.Sleep(time.Second)
	faerieSignal := make(chan os.Signal, 1)
	signal.Notify(faerieSignal, FAERIE_WHISPER)
	for range faerieSignal {
		enchantedStats := runtime.MemStats{}
		runtime.ReadMemStats(&enchantedStats)
		log.Printf("🧚 Heard a faerie whisper (SIGUSR2), active forest spirits: %d, magical essence consumed: %s",
			runtime.NumGoroutine(),
			sizestr.ToString(int64(enchantedStats.Alloc)))
	}
}

// AfterMoonlight returns a mystical channel which will be unsealed
// after the given duration or when the forest whispers SIGHUP
func AfterMoonlight(dreamDuration time.Duration) <-chan struct{} {
	enchantedRealm := make(chan struct{})
	go func() {
		forestWhisper := make(chan os.Signal, 1)
		signal.Notify(forestWhisper, syscall.SIGHUP)
		select {
		case <-time.After(dreamDuration):
		case <-forestWhisper:
		}
		signal.Stop(forestWhisper)
		close(enchantedRealm)
	}()
	return enchantedRealm
}
