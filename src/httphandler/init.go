package httphandler

import (
	"errors"

	"rabbitsky/src/channel"
	"rabbitsky/src/websocket"
)

func Init(c *channel.Channel, ws *websocket.WebSocket, origin, serverPassword string) (*HTTPHandler, error) {
	if c == nil {
		return nil, errors.New("[Error] Init HTTP Handler Failed: Channel is nil!")
	}

	if ws == nil {
		return nil, errors.New("[Error] Init HTTP Handler Failed: WebSocket is nil!")
	}

	if origin == "" {
		return nil, errors.New("[Error] Init HTTP Handler Failed: Origin is Empty!")
	}

	httpHandler := HTTPHandler{
		Channel:        c,
		WebSocket:      ws,
		Origin:         origin,
		ServerPassword: serverPassword,
	}

	return &httpHandler, nil
}
