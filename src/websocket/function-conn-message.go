package websocket

import (
	"time"

	"github.com/gorilla/websocket"
)

func (c *WebSocketConn) DisconnectMessage(str string) {
	c.Mutex.Lock()
	defer c.Mutex.Unlock()

	c.Conn.WriteMessage(websocket.TextMessage, []byte(str))
	c.Conn.WriteMessage(websocket.CloseMessage, nil)
}

func (c *WebSocketConn) SendMessage(str string) {
	c.Mutex.Lock()
	defer c.Mutex.Unlock()

	c.Conn.WriteMessage(websocket.TextMessage, []byte(str))
}

func (c *WebSocketConn) ReadMessage() (int, []byte, error) {
	return c.Conn.ReadMessage()
}

func (c *WebSocketConn) SetReadDeadline(t time.Duration) error {
	return c.Conn.SetReadDeadline(time.Now().Add(t))
}

func (c *WebSocketConn) SetWriteDeadline(t time.Duration) error {
	return c.Conn.SetWriteDeadline(time.Now().Add(t))
}

func (c *WebSocketConn) Close() error {
	return c.Conn.Close()
}
