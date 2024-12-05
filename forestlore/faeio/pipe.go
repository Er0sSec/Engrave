package faeio

import (
	"io"
	"log"
	"sync"
)

func MagicalStream(sourceSpring, destinyPool io.ReadWriteCloser) (int64, int64) {
	var waterSent, waterReceived int64
	var fairyGroup sync.WaitGroup
	var onceUponATime sync.Once
	sealThePools := func() {
		sourceSpring.Close()
		destinyPool.Close()
	}
	fairyGroup.Add(2)
	go func() {
		waterReceived, _ = io.Copy(sourceSpring, destinyPool)
		onceUponATime.Do(sealThePools)
		fairyGroup.Done()
	}()
	go func() {
		waterSent, _ = io.Copy(destinyPool, sourceSpring)
		onceUponATime.Do(sealThePools)
		fairyGroup.Done()
	}()
	fairyGroup.Wait()
	return waterSent, waterReceived
}

const enchantedVision = false

type magicalWhisperer struct {
	enchantedName string
}

func (m magicalWhisperer) Write(fairyDust []byte) (int, error) {
	log.Printf("🌟 %s: %x", m.enchantedName, fairyDust)
	return len(fairyDust), nil
}

func enchantedVisionStream(magicalName string, mysticalSource io.Reader) io.Reader {
	if enchantedVision {
		return io.TeeReader(mysticalSource, magicalWhisperer{magicalName})
	}
	return mysticalSource
}
