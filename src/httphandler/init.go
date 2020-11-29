package httphandler

import (
	"errors"
	"regexp"
	"strconv"
	"strings"

	"rabbitsky/src/channel"
	"rabbitsky/src/websocket"
)

func Init(c *channel.Channel, ws *websocket.WebSocket, origin, serverPassword, positionMin, positionMax string) (*HTTPHandler, error) {
	if c == nil {
		return nil, errors.New("[Error] Init HTTP Handler Failed: Channel is nil!")
	}

	if ws == nil {
		return nil, errors.New("[Error] Init HTTP Handler Failed: WebSocket is nil!")
	}

	if origin == "" {
		return nil, errors.New("[Error] Init HTTP Handler Failed: Origin is Empty!")
	}

	if positionMin == "" {
		return nil, errors.New("[Error] Init HTTP Handler Failed: Position Minimum is Empty!")
	}

	if positionMax == "" {
		return nil, errors.New("[Error] Init HTTP Handler Failed: Position Maximum is Empty!")
	}

	posMinSplit := strings.SplitN(positionMin, ",", 3)
	if len(posMinSplit) != 3 {
		return nil, errors.New("[Error] Init HTTP Handler Failed: Invalid Position Minimum Format!")
	}

	posMaxSplit := strings.SplitN(positionMax, ",", 3)
	if len(posMaxSplit) != 3 {
		return nil, errors.New("[Error] Init HTTP Handler Failed: Invalid Position Maximum Format!")
	}

	posMin := Position{}
	posMax := Position{}

	for i, v := range posMinSplit {
		posInt, err := strconv.Atoi(v)
		if err != nil {
			return nil, errors.New("[Error] Init HTTP Handler Failed: Invalid Position Minimum Format!")
		}

		switch i {
		case 0:
			posMin.X = posInt
			break
		case 1:
			posMin.Y = posInt
			break
		case 2:
			posMin.Z = posInt
			break
		}
	}

	for i, v := range posMaxSplit {
		posInt, err := strconv.Atoi(v)
		if err != nil {
			return nil, errors.New("[Error] Init HTTP Handler Failed: Invalid Position Maximum Format!")
		}

		switch i {
		case 0:
			posMax.X = posInt
			break
		case 1:
			posMax.Y = posInt
			break
		case 2:
			posMax.Z = posInt
			break
		}
	}

	httpHandler := HTTPHandler{
		Channel:        c,
		WebSocket:      ws,
		Origin:         origin,
		ServerPassword: serverPassword,
		MessageRegex:   regexp.MustCompile(`^(H(?P<H>[0-9]{1,3}))?(S(?P<S>[0-9]{1,3}))?(L(?P<L>[0-9]{1,3}))?(X(?P<X>\-?[0-9]+))?(Y(?P<Y>\-?[0-9]+))?(Z(?P<Z>\-?[0-9]+))?(x(?P<x>\-?[0-9]+))?(y(?P<y>\-?[0-9]+))?(z(?P<z>\-?[0-9]+))?(D(?P<D>(0|1)))?(C(?P<C>.*))?$`),

		PositionLimitMin: &posMin,
		PositionLimitMax: &posMax,
	}

	return &httpHandler, nil
}
