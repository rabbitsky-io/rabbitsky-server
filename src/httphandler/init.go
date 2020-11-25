package httphandler

import (
	"errors"
	"regexp"

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
		MessageRegex:   regexp.MustCompile(`^(H(?P<H>[0-9]{1,3}))?(S(?P<S>[0-9]{1,3}))?(L(?P<L>[0-9]{1,3}))?(X(?P<X>\-?[0-9]+))?(Y(?P<Y>\-?[0-9]+))?(Z(?P<Z>\-?[0-9]+))?(x(?P<x>\-?[0-9]+))?(y(?P<y>\-?[0-9]+))?(z(?P<z>\-?[0-9]+))?(D(?P<D>(0|1)))?(C(?P<C>.*))?$`),
	}

	return &httpHandler, nil
}
