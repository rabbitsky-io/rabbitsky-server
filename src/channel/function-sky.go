package channel

import (
	"errors"
)

func (c *Channel) ChangeSkyColor(color string) error {
	if color == "" {
		return errors.New("Empty Color")
	}

	c.SkyColor = color
	return nil
}

func (c *Channel) GetSkyColor() string {
	return c.SkyColor
}
