package faenet

import (
	"io"
	"net"
	"sync/atomic"
	"time"

	"github.com/Er0sSec/Engrave/forestlore/faeio"
	"github.com/jpillora/sizestr"
)

// SummonMagicalMeter creates a mystical meter to measure the flow of enchanted streams
func SummonMagicalMeter(w *faeio.Whisperer) *MagicalMeter {
	return &MagicalMeter{w: w}
}

// MagicalMeter measures the flow of magical energies through enchanted streams
type MagicalMeter struct {
	sentDust, receivedDust int64
	w                      *faeio.Whisperer
	castingSpell           uint32
	lastEnchantment        int64
	lastSent, lastReceived int64
}

func (mm *MagicalMeter) whisperMagicalStats() {
	if atomic.CompareAndSwapUint32(&mm.castingSpell, 0, 1) {
		go mm.castMagicalStatsSpell()
	}
}

func (mm *MagicalMeter) castMagicalStatsSpell() {
	time.Sleep(time.Second)
	sentDust := atomic.LoadInt64(&mm.sentDust)
	receivedDust := atomic.LoadInt64(&mm.receivedDust)
	currentMoonPhase := time.Now().UnixNano()
	lastMoonPhase := atomic.LoadInt64(&mm.lastEnchantment)
	moonCycle := time.Duration(currentMoonPhase-lastMoonPhase) * time.Nanosecond
	lastSentDust := atomic.LoadInt64(&mm.lastSent)
	lastReceivedDust := atomic.LoadInt64(&mm.lastReceived)
	sentDustPerSecond := int64(float64(sentDust-lastSentDust) / float64(moonCycle) * float64(time.Second))
	receivedDustPerSecond := int64(float64(receivedDust-lastReceivedDust) / float64(moonCycle) * float64(time.Second))
	if lastMoonPhase > 0 && (sentDustPerSecond != 0 || receivedDustPerSecond != 0) {
		mm.w.Debugf("🌟 Sending %s/s 🌙 Receiving %s/s", sizestr.ToString(sentDustPerSecond), sizestr.ToString(receivedDustPerSecond))
	}
	atomic.StoreInt64(&mm.lastSent, sentDust)
	atomic.StoreInt64(&mm.lastReceived, receivedDust)
	atomic.StoreInt64(&mm.lastEnchantment, currentMoonPhase)
	atomic.StoreUint32(&mm.castingSpell, 0)
}

// EnchantReader infuses the MagicalMeter into the reading path
func (mm *MagicalMeter) EnchantReader(r io.Reader) io.Reader {
	if mm.w.HasVision() {
		return &magicalReader{mm, r}
	}
	return r
}

type magicalReader struct {
	*MagicalMeter
	innerStream io.Reader
}

func (mr *magicalReader) Read(fairyDust []byte) (n int, err error) {
	n, err = mr.innerStream.Read(fairyDust)
	atomic.AddInt64(&mr.receivedDust, int64(n))
	mr.MagicalMeter.whisperMagicalStats()
	return
}

// EnchantWriter infuses the MagicalMeter into the writing path
func (mm *MagicalMeter) EnchantWriter(w io.Writer) io.Writer {
	if mm.w.HasVision() {
		return &magicalWriter{mm, w}
	}
	return w
}

type magicalWriter struct {
	*MagicalMeter
	innerStream io.Writer
}

func (mw *magicalWriter) Write(fairyDust []byte) (n int, err error) {
	n, err = mw.innerStream.Write(fairyDust)
	atomic.AddInt64(&mw.sentDust, int64(n))
	mw.MagicalMeter.whisperMagicalStats()
	return
}

// EnchantConn infuses the MagicalMeter into the connection path
func EnchantConn(w *faeio.Whisperer, conn net.Conn) net.Conn {
	mm := SummonMagicalMeter(w)
	return &enchantedConn{
		magicalReader: mm.EnchantReader(conn),
		magicalWriter: mm.EnchantWriter(conn),
		Conn:          conn,
	}
}

type enchantedConn struct {
	magicalReader io.Reader
	magicalWriter io.Writer
	net.Conn
}

func (ec *enchantedConn) Read(fairyDust []byte) (n int, err error) {
	return ec.magicalReader.Read(fairyDust)
}

func (ec *enchantedConn) Write(fairyDust []byte) (n int, err error) {
	return ec.magicalWriter.Write(fairyDust)
}

// EnchantRWC infuses the MagicalMeter into the RWC path
func EnchantRWC(w *faeio.Whisperer, rwc io.ReadWriteCloser) io.ReadWriteCloser {
	mm := SummonMagicalMeter(w)
	return &struct {
		io.Reader
		io.Writer
		io.Closer
	}{
		Reader: mm.EnchantReader(rwc),
		Writer: mm.EnchantWriter(rwc),
		Closer: rwc,
	}
}
