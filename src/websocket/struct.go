package websocket

import (
	"sync"

	gws "github.com/gorilla/websocket"
)

type WebSocket struct {
	Upgrader *gws.Upgrader
	Origin   string
}

type WebSocketConn struct {
	Conn  *gws.Conn
	Mutex sync.Mutex
}
