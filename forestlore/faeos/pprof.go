//go:build pprof
// +build pprof

package faeOS

import (
	"log"
	"net/http"
	_ "net/http/pprof" // summon the mystical profiler spirits
)

func init() {
	go func() {
		log.Fatal(http.ListenAndServe("localhost:6060", nil))
	}()
	log.Printf("🌟 [Enchanted Profiler] whispering secrets on magical portal 6060")
}
