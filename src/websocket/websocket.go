package websocket

import (
	"errors"
	"net/http"

	gws "github.com/gorilla/websocket"
)

func Init(origin string) (*WebSocket, error) {
	if origin == "" {
		return nil, errors.New("[Error] Init WebSocket Failed: Origin is empty!")
	}

	ws := WebSocket{}

	ws.Origin = origin
	ws.Upgrader = &gws.Upgrader{
		EnableCompression: true,
		ReadBufferSize:    16384,
		WriteBufferSize:   16384,
	}

	ws.Upgrader.CheckOrigin = ws.checkOrigin

	return &ws, nil
}

func (ws *WebSocket) checkOrigin(r *http.Request) bool {
	if r.Header.Get("Origin") == ws.Origin {
		return true
	}

	return false
}

func (ws *WebSocket) Upgrade(w http.ResponseWriter, r *http.Request) (*WebSocketConn, error) {
	conn, err := ws.Upgrader.Upgrade(w, r, nil)
	if err != nil {
		return nil, err
	}

	wsConn := WebSocketConn{
		Conn: conn,
	}

	return &wsConn, nil
}
