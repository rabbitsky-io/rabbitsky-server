package player_test

import (
	rsPlayer "rabbitsky/src/player"
	"testing"
)

func TestCreate(t *testing.T) {
	playerObj := rsPlayer.Create()
	if playerObj == nil {
		t.Fail()
		return
	}
	return
}
