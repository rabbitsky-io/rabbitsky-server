package channel

import (
	"errors"

	"rabbitsky/src/player"

	"github.com/catinello/base62"
)

func (c *Channel) GetMaxPlayers() int {
	return c.MaxPlayers
}

func (c *Channel) GetCurrentPlayers() int {
	return c.Players.Count()
}

func (c *Channel) CreatePlayer() (*player.Player, error) {
	if c.GetCurrentPlayers() >= c.GetMaxPlayers() {
		return nil, errors.New("Server is full")
	}

	c.LastID++
	id := c.LastID

	player := player.Create()
	player.ID = base62.Encode(id)

	c.Players.Set(player.ID, player)

	return player, nil
}

func (c *Channel) RemovePlayer(id string) error {
	if id == "" {
		return errors.New("ID is empty")
	}

	c.Players.Remove(id)

	return nil
}
