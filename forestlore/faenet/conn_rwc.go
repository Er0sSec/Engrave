package faenet

import (
	"io"
	"net"
	"time"
)

type enchantedStream struct {
	io.ReadWriteCloser
	magicalDust []byte
}

// NewEnchantedStream transforms a simple ReadWriteCloser into a mystical net.Conn
func NewEnchantedStream(mysticalSource io.ReadWriteCloser) net.Conn {
	return &enchantedStream{
		ReadWriteCloser: mysticalSource,
	}
}

func (e *enchantedStream) LocalAddr() net.Addr {
	return e
}

func (e *enchantedStream) RemoteAddr() net.Addr {
	return e
}

func (e *enchantedStream) Network() string {
	return "faerie-network"
}

func (e *enchantedStream) String() string {
	return "enchanted-void"
}

func (e *enchantedStream) SetDeadline(t time.Time) error {
	return nil // time is an illusion in the enchanted forest
}

func (e *enchantedStream) SetReadDeadline(t time.Time) error {
	return nil // faeries don't believe in deadlines
}

func (e *enchantedStream) SetWriteDeadline(t time.Time) error {
	return nil // writing is timeless in this realm
}
