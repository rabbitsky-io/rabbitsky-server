package httphandler

import (
	"rabbitsky/src/channel"
	"rabbitsky/src/websocket"
)

type HTTPHandler struct {
	Channel        *channel.Channel
	WebSocket      *websocket.WebSocket
	Origin         string
	ServerPassword string
}

type ChannelPlayersJSON struct {
	MaxPlayers int `json:"max_players"`
	Players    int `json:"players"`
}
