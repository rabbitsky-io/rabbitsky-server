package player

import (
	"rabbitsky/src/websocket"
)

type Player struct {
	ID         string
	WebSocket  *websocket.WebSocketConn
	Ready      bool
	ColorR     int
	ColorG     int
	ColorB     int
	PosX       float64
	PosY       float64
	PosZ       float64
	LookX      float64
	LookY      float64
	LookZ      float64
	IsDuck     bool
	Chat       string
	UpdateSent bool
}
