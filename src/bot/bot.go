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
			player.PosX = float64(2000 - rand.Intn(4000))
			player.PosY = float64(10)
			player.PosZ = float64(1500 - rand.Intn(3000))
			player.LookX = 0
			player.LookY = 1000
			player.LookZ = 3000
			player.IsDuck = false
			player.Ready = true

			playerBot = append(playerBot, player)
		}
	}

	tickToMS := time.Duration(1000 / tick)
	ticker := time.NewTicker(tickToMS * time.Millisecond)

	for ; true; <-ticker.C {
		for i := 0; i < maxBot; i++ {
			playerBot[i].PosX = playerBot[i].PosX + float64(25-rand.Intn(50))
			if playerBot[i].PosX > 2000 {
				playerBot[i].PosX = 2000
			}
			if playerBot[i].PosX < -2000 {
				playerBot[i].PosX = -2000
			}

			playerBot[i].PosY = float64(rand.Intn(5) + 10)
			playerBot[i].PosZ = playerBot[i].PosZ + float64(25-rand.Intn(50))
			if playerBot[i].PosZ > 1500 {
				playerBot[i].PosZ = 1500
			}
			if playerBot[i].PosZ < -1500 {
				playerBot[i].PosZ = -1500
			}

			if rand.Intn(25) == 0 {
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

			sendText := fmt.Sprintf("%s,%s,%.2f,%.2f,%.2f,%.2f,%.2f,%.2f,%d,%s",
				rsHTTPHandler.SEND_PLAYER_UPDATE,
				playerBot[i].ID,
				playerBot[i].PosX,
				playerBot[i].PosY,
				playerBot[i].PosZ,
				playerBot[i].LookX,
				playerBot[i].LookY,
				playerBot[i].LookZ,
				isDuck,
				playerBot[i].Chat,
			)

			c.AddBroadcastMessage(sendText)
		}
	}
}
