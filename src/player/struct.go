package player

import (
	"rabbitsky/src/websocket"
)

type Player struct {
	ID         string
	WebSocket  *websocket.WebSocketConn
	IsAdmin    bool
	Ready      bool
	Size       int
	ColorH     int
	ColorS     int
	ColorL     int
	PosX       int
	PosY       int
	PosZ       int
	LookX      int
	LookY      int
	LookZ      int
	IsDuck     bool
	Chat       string
	UpdateSent bool
	AdminTry   int
}
