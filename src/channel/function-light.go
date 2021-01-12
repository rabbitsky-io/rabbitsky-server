package channel

import (
	"errors"
)

func (c *Channel) ChangeLightState(state string) error {
	if state == "" {
		return errors.New("Empty Light State")
	}

	c.LightState = state
	return nil
}

func (c *Channel) GetLightState() string {
	return c.LightState
}
