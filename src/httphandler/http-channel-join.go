package httphandler

import (
	"errors"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"strings"
	"time"

	// "time"

	rsPlayer "rabbitsky/src/player"
)

func (h *HTTPHandler) ChannelJoin(w http.ResponseWriter, r *http.Request) {
	conn, connErr := h.WebSocket.Upgrade(w, r)
	if connErr != nil {
		log.Println("Err: Cannot Upgrade -", connErr)
		return
	}

	defer conn.Close()

	player, err := h.Channel.CreatePlayer()
	if err != nil {
		disconnectMessage := fmt.Sprintf("%s,%s", SEND_DISCONNECT, err.Error())
		conn.DisconnectMessage(disconnectMessage)
		return
	}

	pingTime := 5 * time.Second
	pingTicker := time.NewTicker(pingTime)
	pingDone := make(chan bool)

	go func() {
		for {
			select {
			case <-pingDone:
				return
			case <-pingTicker.C:
				player.SendMessage(SEND_PING)
			}
		}
	}()

	defer func() {
		pingTicker.Stop()
		pingDone <- true

		h.Channel.RemovePlayer(player.ID)

		broadcastMessage := fmt.Sprintf("%s,%s", SEND_PLAYER_DISCONNECT, player.ID)
		h.Channel.AddBroadcastMessage(broadcastMessage)
	}()

	player.WebSocket = conn

	sendInitMessage := fmt.Sprintf("%s,%s", SEND_PLAYER_ID, player.ID)
	player.SendMessage(sendInitMessage)

	wsTimeout := 15 * time.Second

	for {
		message, err := player.ReadMessage()
		if err != nil {
			errDisconnect := fmt.Sprintf("%s,%s", SEND_DISCONNECT, err.Error())
			player.DisconnectMessage(errDisconnect)
			break
		}

		player.SetReadDeadline(wsTimeout)
		player.SetWriteDeadline(wsTimeout)

		if message == "" {
			continue
		}

		messageSplit := strings.SplitN(message, ",", 2)

		if len(messageSplit) <= 0 {
			continue
		}

		if messageSplit[0] == RECEIVE_PING {
			continue
		}

		if messageSplit[0] == RECEIVE_PLAYER_INIT {
			messageSplitFormat := strings.SplitN(message, ",", 11)

			err = h.MessageInitPlayer(player, messageSplitFormat)
			if err != nil {
				errDisconnect := fmt.Sprintf("%s,%s", SEND_DISCONNECT, err.Error())
				player.DisconnectMessage(errDisconnect)
				break
			}

			broadcastMessage := fmt.Sprintf("%s,%s,%s", SEND_PLAYER_INIT, player.ID, strings.Join(messageSplitFormat[1:], ","))
			h.Channel.AddBroadcastMessage(broadcastMessage)

			player.Ready = true

			// Continue so no more checking to bottom
			continue
		}

		if messageSplit[0] == RECEIVE_PLAYER_UPDATE {
			if !player.Ready {
				continue
			}

			// Do not send another update until next tick
			if player.UpdateSent {
				continue
			}

			messageSplitFormat := strings.SplitN(message, ",", 9)

			err = h.MessageUpdatePlayer(player, messageSplitFormat)
			if err != nil {
				continue
			}

			/* Chat parse command */
			err = h.ParseCommand(player, messageSplitFormat[:])
			if err != nil {
				continue
			}
			broadcastMessage := fmt.Sprintf("%s,%s,%s", SEND_PLAYER_UPDATE, player.ID, strings.Join(messageSplitFormat[1:], ","))
			h.Channel.AddBroadcastMessage(broadcastMessage)

			player.UpdateSent = true

			// Continue so no more checking to bottom
			continue
		}
	}
}

func (h *HTTPHandler) ParseCommand(player *rsPlayer.Player, data []string) error {
	if player == nil {
		return errors.New("Player is nil")
	}

	if len(data) != 9 {
		return errors.New("Invalid Data")
	}

	chat := data[8]

	if chat == "" {
		return nil
	}

	if chat[0] == '/' {
		data[8] = ""

		if h.ServerPassword != "" {
			command := strings.Split(chat, " ")
			if len(command) > 1 {
				switch command[0] {
				case "/admin":
					if command[1] == h.ServerPassword {
						player.IsAdmin = true
					}
					break
				case "/sky":
					if player.IsAdmin && command[1] != "" {
						h.Channel.ChangeSkyColor(command[1])

						broadcastMessage := fmt.Sprintf("%s,%s", SEND_SKY_COLOR, command[1])
						h.Channel.AddBroadcastMessage(broadcastMessage)
					}
					break
				}
			}
		}
	}

	return nil
}

func (h *HTTPHandler) MessageInitPlayer(player *rsPlayer.Player, data []string) error {
	if player == nil {
		return errors.New("Player is nil")
	}

	if len(data) != 11 {
		return errors.New("Invalid Data")
	}

	if data[1] == "" || data[2] == "" || data[3] == "" || data[4] == "" || data[5] == "" || data[6] == "" || data[7] == "" || data[8] == "" || data[9] == "" {
		return errors.New("Found empty data on init")
	}

	if colorR, err := strconv.Atoi(data[1]); err == nil {
		player.ColorR = colorR
	}

	if colorG, err := strconv.Atoi(data[2]); err == nil {
		player.ColorG = colorG
	}

	if colorB, err := strconv.Atoi(data[3]); err == nil {
		player.ColorB = colorB
	}

	if posX, err := strconv.ParseFloat(data[4], 64); err == nil {
		player.PosX = posX
	}

	if posY, err := strconv.ParseFloat(data[5], 64); err == nil {
		player.PosY = posY
	}

	if posZ, err := strconv.ParseFloat(data[6], 64); err == nil {
		player.PosZ = posZ
	}

	if LookX, err := strconv.ParseFloat(data[7], 64); err == nil {
		player.LookX = LookX
	}

	if LookY, err := strconv.ParseFloat(data[8], 64); err == nil {
		player.LookY = LookY
	}

	if LookZ, err := strconv.ParseFloat(data[9], 64); err == nil {
		player.LookZ = LookZ
	}

	if data[10] != "" {
		if data[10] == "1" {
			player.IsDuck = true
		} else {
			player.IsDuck = false
		}
	}

	h.SendInit(player)

	return nil
}

func (h *HTTPHandler) MessageUpdatePlayer(player *rsPlayer.Player, data []string) error {
	if player == nil {
		return errors.New("Player is nil")
	}

	if len(data) != 9 {
		return errors.New("Invalid Data")
	}

	if data[1] != "" {
		if posX, err := strconv.ParseFloat(data[1], 64); err == nil {
			player.PosX = posX
		}
	}

	if data[2] != "" {
		if posY, err := strconv.ParseFloat(data[2], 64); err == nil {
			player.PosY = posY
		}
	}

	if data[3] != "" {
		if posZ, err := strconv.ParseFloat(data[3], 64); err == nil {
			player.PosZ = posZ
		}
	}

	if data[4] != "" {
		if LookX, err := strconv.ParseFloat(data[4], 64); err == nil {
			player.LookX = LookX
		}
	}

	if data[5] != "" {
		if LookY, err := strconv.ParseFloat(data[5], 64); err == nil {
			player.LookY = LookY
		}
	}

	if data[6] != "" {
		if LookZ, err := strconv.ParseFloat(data[6], 64); err == nil {
			player.LookZ = LookZ
		}
	}

	if data[7] != "" {
		if data[7] == "1" {
			player.IsDuck = true
		} else {
			player.IsDuck = false
		}
	}

	/* We ignore data chat due to bug, do not store! */

	return nil
}

func (h *HTTPHandler) SendInit(player *rsPlayer.Player) error {
	if player == nil {
		return errors.New("Player is nil")
	}

	var str strings.Builder

	for p := range h.Channel.Players.Iter() {
		if p.Val != nil {
			playerObj := p.Val.(*rsPlayer.Player)

			if playerObj.ID == player.ID {
				continue
			}

			isDuck := 0

			if playerObj.IsDuck {
				isDuck = 1
			}

			sendText := fmt.Sprintf("1,%s,%d,%d,%d,%.2f,%.2f,%.2f,%.2f,%.2f,%.2f,%d",
				playerObj.ID,
				playerObj.ColorR,
				playerObj.ColorG,
				playerObj.ColorB,
				playerObj.PosX,
				playerObj.PosY,
				playerObj.PosZ,
				playerObj.LookX,
				playerObj.LookY,
				playerObj.LookZ,
				isDuck,
			)

			if str.Len() > 0 {
				str.WriteString("\n")
			}

			str.WriteString(sendText)
		}
	}

	skyColor := h.Channel.GetSkyColor()
	if skyColor != "" {
		if str.Len() > 0 {
			str.WriteString("\n")
		}

		str.WriteString(fmt.Sprintf("%s,%s", SEND_SKY_COLOR, skyColor))
	}

	player.SendMessage(str.String())

	return nil
}
