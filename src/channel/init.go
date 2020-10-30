package channel

import (
	"errors"
	"log"
	"time"

	cmap "github.com/orcaman/concurrent-map"
)

func Init(maxPlayers, tick int) (*Channel, error) {
	if tick <= 0 {
		return nil, errors.New("[Error] Init Channel Failed: Server tick must be more than zero.")
	}

	if tick > 60 {
		return nil, errors.New("[Error] Init Channel Failed: Server tick must less than equal 60. It's dangerous!")
	}

	// Print Warning!
	if tick > 30 {
		log.Println("[Warning] Init Channel: Server tick more than 30 is a bit too high. Please consider decreasing it, saving bandwidth both users and server.")
	}

	if maxPlayers > 250 {
		log.Println("[Warning] Init Channel: Maximum Players more than 250 is too high. Please consider both the users browser capability and server bandwidth.")
	}

	tickToMS := time.Duration(1000 / tick)
	ticker := time.NewTicker(tickToMS * time.Millisecond)

	channel := Channel{
		ServerTick: tick,
		MaxPlayers: maxPlayers,
		LastID:     0,
		Players:    cmap.New(),
		Ticker:     ticker,
	}

	go channel.TickHandler()

	return &channel, nil
}
