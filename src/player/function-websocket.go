package player

import (
	"errors"
	"strings"
	"time"
)

func (p *Player) DisconnectMessage(str string) {
	if p.WebSocket == nil {
		return
	}

	p.WebSocket.DisconnectMessage(str)
}

func (p *Player) SendMessage(str string) {
	if p.WebSocket == nil {
		return
	}

	p.WebSocket.SendMessage(str)
}

func (p *Player) ReadMessage() (string, error) {
	if p.WebSocket == nil {
		return "", errors.New("No Connection...")
	}

	_, msg, err := p.WebSocket.ReadMessage()
	if err != nil {
		return "", err
	}

	var msgStr strings.Builder
	_, err = msgStr.Write(msg)
	if err != nil {
		return "", err
	}

	return strings.TrimSpace(msgStr.String()), nil
}

func (p *Player) SetReadDeadline(t time.Duration) error {
	return p.WebSocket.SetReadDeadline(t)
}

func (p *Player) SetWriteDeadline(t time.Duration) error {
	return p.WebSocket.SetWriteDeadline(t)
}
