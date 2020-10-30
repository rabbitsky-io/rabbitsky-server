package channel

import (
	rsPlayer "rabbitsky/src/player"
)

func (c *Channel) AddBroadcastMessage(str string) {
	c.BroadcastMutex.Lock()
	defer c.BroadcastMutex.Unlock()

	if c.BroadcastMessage.Len() > 0 {
		c.BroadcastMessage.WriteString("\n")
	}

	c.BroadcastMessage.WriteString(str)
}

func (c *Channel) GetBroadcastMessage() string {
	c.BroadcastMutex.Lock()
	defer c.BroadcastMutex.Unlock()

	msg := c.BroadcastMessage.String()
	c.BroadcastMessage.Reset()

	return msg
}

func (c *Channel) BroadcastUpdate() {
	str := c.GetBroadcastMessage()
	if c.Players.Count() > 0 {
		for p := range c.Players.Iter() {
			if p.Val != nil {
				player := p.Val.(*rsPlayer.Player)
				if player.Ready {
					player.UpdateSent = false
					player.SendMessage(str)
				}
			}
		}
	}
}
