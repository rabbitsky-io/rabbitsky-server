package channel

import (
	"errors"
	"math/rand"

	rsPlayer "rabbitsky/src/player"
)

func (c *Channel) CreateBot(player *rsPlayer.Player) (*rsPlayer.Player, error) {
	if player == nil {
		return nil, errors.New("Player must not empty")
	}

	p, err := c.CreatePlayer()
	if err != nil {
		return nil, err
	}

	c.Bots = append(c.Bots, p)

	p.ColorH = rand.Intn(360)
	p.ColorS = rand.Intn(100)
	p.ColorL = rand.Intn(70) + 15
	p.PosX = player.PosX
	p.PosY = player.PosY
	p.PosZ = player.PosZ
	p.LookX = player.LookX
	p.LookY = player.LookY
	p.LookZ = player.LookZ

	return p, nil
}

func (c *Channel) GetBots() ([]*rsPlayer.Player, error) {
	return c.Bots, nil
}

func (c *Channel) RemoveBots() error {
	for _, v := range c.Bots {
		c.Players.Remove(v.ID)
	}

	c.Bots = nil

	return nil
}
