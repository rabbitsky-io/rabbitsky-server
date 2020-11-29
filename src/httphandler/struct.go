package httphandler

import (
	"regexp"

	"rabbitsky/src/channel"
	"rabbitsky/src/websocket"
)

type HTTPHandler struct {
	Channel        *channel.Channel
	WebSocket      *websocket.WebSocket
	Origin         string
	ServerPassword string
	MessageRegex   *regexp.Regexp

	PositionLimitMin *Position
	PositionLimitMax *Position
}

type Position struct {
	X int
	Y int
	Z int
}

type ChannelPlayersJSON struct {
	MaxPlayers int `json:"max_players"`
	Players    int `json:"players"`
}
