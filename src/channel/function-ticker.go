package channel

func (c *Channel) TickHandler() {
	for ; true; <-c.Ticker.C {
		c.BroadcastUpdate()
	}
}
