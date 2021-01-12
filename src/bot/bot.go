package bot

import (
	"fmt"
	"math/rand"
	"time"

	rsChannel "rabbitsky/src/channel"
	rsHTTPHandler "rabbitsky/src/httphandler"
	rsPlayer "rabbitsky/src/player"
)

func AddBot(c *rsChannel.Channel, maxBot, tick int) {
	if c == nil {
		return
	}

	if maxBot <= 0 {
		return
	}

	var playerBot []*rsPlayer.Player

	for i := 0; i < maxBot; i++ {
		player, err := c.CreatePlayer()
		if err == nil {
			player.ColorH = rand.Intn(360)
			player.ColorS = rand.Intn(100)
			player.ColorL = rand.Intn(70) + 15
			player.PosX = rand.Intn(4000)
			player.PosY = 10
			player.PosZ = rand.Intn(3000)
			player.LookX = 2100
			player.LookY = 1000
			player.LookZ = 6000
			player.IsDuck = false
			player.Ready = true

			playerBot = append(playerBot, player)
		}
	}

	tickToMS := time.Duration(1000 / tick)
	ticker := time.NewTicker(tickToMS * time.Millisecond)

	for ; true; <-ticker.C {
		for i := 0; i < maxBot; i++ {
			playerBot[i].PosX = playerBot[i].PosX + (10 - rand.Intn(20))
			if playerBot[i].PosX > 4000 {
				playerBot[i].PosX = 4000
			}
			if playerBot[i].PosX < 0 {
				playerBot[i].PosX = 0
			}

			playerBot[i].PosY = rand.Intn(5) + 10
			playerBot[i].PosZ = playerBot[i].PosZ + (10 - rand.Intn(20))
			if playerBot[i].PosZ > 3000 {
				playerBot[i].PosZ = 3000
			}
			if playerBot[i].PosZ < 0 {
				playerBot[i].PosZ = 0
			}

			if rand.Intn(100) == 0 {
				playerBot[i].Chat = fmt.Sprintf("ye - %d", rand.Intn(50))
			} else {
				playerBot[i].Chat = ""
			}

			isDuck := 0
			playerBot[i].IsDuck = false

			if rand.Intn(10) < 3 {
				isDuck = 1
				playerBot[i].IsDuck = true
			}

			sendText := fmt.Sprintf("%s%s=X%dY%dZ%dD%dC%s",
				rsHTTPHandler.SEND_PLAYER_UPDATE,
				playerBot[i].ID,
				playerBot[i].PosX,
				playerBot[i].PosY,
				playerBot[i].PosZ,
				isDuck,
				playerBot[i].Chat,
			)

			c.AddBroadcastMessage(sendText)
		}
	}
}
