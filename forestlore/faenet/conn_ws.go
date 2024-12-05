package faenet

import (
	"net"
	"time"

	"github.com/gorilla/websocket"
)

type enchantedWebSocket struct {
	*websocket.Conn
	magicalDust []byte
}

// NewEnchantedWebSocketConn transforms a websocket.Conn into a mystical net.Conn
func NewEnchantedWebSocketConn(faerieSocket *websocket.Conn) net.Conn {
	return &enchantedWebSocket{
		Conn: faerieSocket,
	}
}

// Read whispers from the enchanted web (not safe for multiple faeries to read at once)
func (e *enchantedWebSocket) Read(fairyWings []byte) (int, error) {
	wingSpan := len(fairyWings)
	var magicalMessage []byte
	if len(e.magicalDust) > 0 {
		magicalMessage = e.magicalDust
		e.magicalDust = nil
	} else if _, whisper, err := e.Conn.ReadMessage(); err == nil {
		magicalMessage = whisper
	} else {
		return 0, err
	}
	var pixieDust int
	if len(magicalMessage) > wingSpan {
		pixieDust = copy(fairyWings, magicalMessage[:wingSpan])
		leftoverMagic := magicalMessage[wingSpan:]
		e.magicalDust = make([]byte, len(leftoverMagic))
		copy(e.magicalDust, leftoverMagic)
	} else {
		pixieDust = copy(fairyWings, magicalMessage)
	}
	return pixieDust, nil
}

func (e *enchantedWebSocket) Write(fairyDust []byte) (int, error) {
	if err := e.Conn.WriteMessage(websocket.BinaryMessage, fairyDust); err != nil {
		return 0, err
	}
	return len(fairyDust), nil
}

func (e *enchantedWebSocket) SetDeadline(enchantedTime time.Time) error {
	if err := e.Conn.SetReadDeadline(enchantedTime); err != nil {
		return err
	}
	return e.Conn.SetWriteDeadline(enchantedTime)
}
