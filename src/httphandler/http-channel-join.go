package httphandler

import (
	"errors"
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"time"

	rsPlayer "rabbitsky/src/player"
)

func (h *HTTPHandler) ChannelJoin(w http.ResponseWriter, r *http.Request) {
	conn, connErr := h.WebSocket.Upgrade(w, r)
	if connErr != nil {
		return
	}

	defer conn.Close()

	player, err := h.Channel.CreatePlayer()
	if err != nil {
		disconnectMessage := fmt.Sprintf("%s%s", SEND_DISCONNECT, err.Error())
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

		broadcastMessage := fmt.Sprintf("%s%s", SEND_PLAYER_DISCONNECT, player.ID)
		h.Channel.AddBroadcastMessage(broadcastMessage)
	}()

	player.WebSocket = conn

	sendInitMessage := fmt.Sprintf("%s%s", SEND_PLAYER_ID, player.ID)
	player.SendMessage(sendInitMessage)

	wsTimeout := 15 * time.Second

	for {
		message, err := player.ReadMessage()
		if err != nil {
			errDisconnect := fmt.Sprintf("%s%s", SEND_DISCONNECT, err.Error())
			player.DisconnectMessage(errDisconnect)
			break
		}

		player.SetReadDeadline(wsTimeout)
		player.SetWriteDeadline(wsTimeout)

		if message == "" {
			continue
		}

		messageCode := string(message[0])

		if messageCode == RECEIVE_PING {
			continue
		}

		if messageCode == RECEIVE_PLAYER_INIT {
			if player.Ready {
				continue
			}

			if len(message) < 2 {
				errDisconnect := fmt.Sprintf("%s%s", SEND_DISCONNECT, "Invalid Data")
				player.DisconnectMessage(errDisconnect)
				break
			}

			messageData := message[1:len(message)]
			if messageData == "" {
				errDisconnect := fmt.Sprintf("%s%s", SEND_DISCONNECT, "Invalid Data")
				player.DisconnectMessage(errDisconnect)
				break
			}

			parsedData, err := h.ParseMessageData(messageData)
			if err != nil {
				errDisconnect := fmt.Sprintf("%s%s", SEND_DISCONNECT, err.Error())
				player.DisconnectMessage(errDisconnect)
				break
			}

			err = h.MessageInitPlayer(player, parsedData)
			if err != nil {
				errDisconnect := fmt.Sprintf("%s%s", SEND_DISCONNECT, err.Error())
				player.DisconnectMessage(errDisconnect)
				break
			}

			unparsedData := h.UnparseMessageData(parsedData)
			if unparsedData == "" {
				continue
			}

			broadcastMessage := fmt.Sprintf("%s%s=%s", SEND_PLAYER_INIT, player.ID, unparsedData)
			h.Channel.AddBroadcastMessage(broadcastMessage)

			player.Ready = true

			// Continue so no more checking to bottom
			continue
		}

		if messageCode == RECEIVE_PLAYER_UPDATE {
			if !player.Ready {
				continue
			}

			// Do not send another update until next tick
			if player.UpdateSent {
				continue
			}

			if len(message) < 2 {
				continue
			}

			messageData := message[1:len(message)]
			if messageData == "" {
				continue
			}

			parsedData, err := h.ParseMessageData(messageData)
			if err != nil {
				errDisconnect := fmt.Sprintf("%s%s", SEND_DISCONNECT, err.Error())
				player.DisconnectMessage(errDisconnect)
				break
			}

			err = h.MessageUpdatePlayer(player, parsedData)
			if err != nil {
				continue
			}

			/* Chat parse command */
			h.ParseCommand(player, parsedData)

			unparsedData := h.UnparseMessageData(parsedData)
			if unparsedData == "" {
				continue
			}

			broadcastMessage := fmt.Sprintf("%s%s=%s", SEND_PLAYER_UPDATE, player.ID, unparsedData)
			h.Channel.AddBroadcastMessage(broadcastMessage)

			player.UpdateSent = true

			// Continue so no more checking to bottom
			continue
		}
	}
}

func (h *HTTPHandler) ParseCommand(player *rsPlayer.Player, data map[string]string) error {
	if player == nil {
		return errors.New("Player is nil")
	}

	chat := data["C"]
	if chat == "" {
		return errors.New("Invalid Data")
	}

	if chat[0] != '/' {
		return errors.New("Not a command")
	}

	data["C"] = "" // Delete Map

	if h.ServerPassword == "" {
		return errors.New("Server Password Empty")
	}

	command := strings.Split(chat, " ")
	if len(command) > 0 {
		switch command[0] {
		case "/admin":
			if command[1] == h.ServerPassword {
				player.IsAdmin = true
			}

			break
		case "/sky":
			if player.IsAdmin && command[1] != "" {
				skySend := fmt.Sprintf("%s%s", SEND_SKY_COLOR_STANDARD, command[1])
				h.Channel.ChangeSkyColor(skySend)

				broadcastMessage := fmt.Sprintf("%s%s", SEND_SKY_COLOR, skySend)
				h.Channel.AddBroadcastMessage(broadcastMessage)
			}

			break
		case "/skyflash":
			if player.IsAdmin && len(command) > 2 {
				_, err := strconv.Atoi(command[1])
				if err != nil {
					break
				}

				commandSaveArr := []string{}

				_, err = strconv.Atoi(command[2])
				if err == nil {
					commandSaveArr = append(commandSaveArr, command[1:]...)
				} else {
					commandSaveArr = append(commandSaveArr, command[1])
					commandSaveArr = append(commandSaveArr, command[1:]...)
				}

				skySend := fmt.Sprintf("%s%s", SEND_SKY_COLOR_FLASH, strings.Join(commandSaveArr, ","))
				h.Channel.ChangeSkyColor(skySend)

				broadcastMessage := fmt.Sprintf("%s%s", SEND_SKY_COLOR, skySend)
				h.Channel.AddBroadcastMessage(broadcastMessage)
			}

			break
		case "/botadd":
			if player.IsAdmin {
				bot, err := h.Channel.CreateBot(player)
				if err != nil {
					break
				}

				broadcastMessage := fmt.Sprintf("%s%s=H%dS%dL%dX%dY%dZ%dx%dy%dz%dD%d",
					SEND_PLAYER_INIT,
					bot.ID,
					bot.ColorH,
					bot.ColorS,
					bot.ColorL,
					bot.PosX,
					bot.PosY,
					bot.PosZ,
					bot.LookX,
					bot.LookY,
					bot.LookZ,
					0,
				)

				h.Channel.AddBroadcastMessage(broadcastMessage)
			}

			break
		case "/botremove":
			if player.IsAdmin {
				bots, err := h.Channel.GetBots()
				if err != nil {
					break
				}

				err = h.Channel.RemoveBots()

				for _, bot := range bots {
					broadcastMessage := fmt.Sprintf("%s%s", SEND_PLAYER_DISCONNECT, bot.ID)
					h.Channel.AddBroadcastMessage(broadcastMessage)
				}
			}

			break
		}
	}

	return nil
}

func (h *HTTPHandler) MessageInitPlayer(player *rsPlayer.Player, data map[string]string) error {
	if player == nil {
		return errors.New("Player is nil")
	}

	for k, v := range data {
		if k == "C" {
			continue
		}

		if v == "" {
			return errors.New("Found Empty Data")
		}
	}

	if colorH, err := strconv.Atoi(data["H"]); err == nil {
		player.ColorH = colorH
	}

	if colorS, err := strconv.Atoi(data["S"]); err == nil {
		player.ColorS = colorS
	}

	if colorL, err := strconv.Atoi(data["L"]); err == nil {
		player.ColorL = colorL
	}

	if posX, err := strconv.Atoi(data["X"]); err == nil {
		player.PosX = posX
	}

	if posY, err := strconv.Atoi(data["Y"]); err == nil {
		player.PosY = posY
	}

	if posZ, err := strconv.Atoi(data["Z"]); err == nil {
		player.PosZ = posZ
	}

	if LookX, err := strconv.Atoi(data["x"]); err == nil {
		player.LookX = LookX
	}

	if LookY, err := strconv.Atoi(data["y"]); err == nil {
		player.LookY = LookY
	}

	if LookZ, err := strconv.Atoi(data["z"]); err == nil {
		player.LookZ = LookZ
	}

	if data["D"] != "" {
		if data["D"] == "1" {
			player.IsDuck = true
		} else {
			player.IsDuck = false
		}
	}

	h.SendInit(player)

	return nil
}

func (h *HTTPHandler) MessageUpdatePlayer(player *rsPlayer.Player, data map[string]string) error {
	if player == nil {
		return errors.New("Player is nil")
	}

	if data["X"] != "" {
		if posX, err := strconv.Atoi(data["X"]); err == nil {
			player.PosX = posX
		}
	}

	if data["Y"] != "" {
		if posY, err := strconv.Atoi(data["Y"]); err == nil {
			player.PosY = posY
		}
	}

	if data["Z"] != "" {
		if posZ, err := strconv.Atoi(data["Z"]); err == nil {
			player.PosZ = posZ
		}
	}

	if data["x"] != "" {
		if LookX, err := strconv.Atoi(data["x"]); err == nil {
			player.LookX = LookX
		}
	}

	if data["y"] != "" {
		if LookY, err := strconv.Atoi(data["y"]); err == nil {
			player.LookY = LookY
		}
	}

	if data["z"] != "" {
		if LookZ, err := strconv.Atoi(data["z"]); err == nil {
			player.LookZ = LookZ
		}
	}

	if data["D"] != "" {
		if data["D"] == "1" {
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

			sendText := fmt.Sprintf("%s%s=H%dS%dL%dX%dY%dZ%dx%dy%dz%dD%d",
				SEND_PLAYER_INIT,
				playerObj.ID,
				playerObj.ColorH,
				playerObj.ColorS,
				playerObj.ColorL,
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

		str.WriteString(fmt.Sprintf("%s%s", SEND_SKY_COLOR, skyColor))
	}

	player.SendMessage(str.String())

	return nil
}

func (h *HTTPHandler) ParseMessageData(str string) (map[string]string, error) {
	if str == "" {
		return nil, errors.New("String is empty")
	}

	regexMatch := h.MessageRegex.FindStringSubmatch(str)
	if len(regexMatch) == 0 {
		return nil, errors.New("Data not match")
	}

	playerData := map[string]string{}
	for i, name := range h.MessageRegex.SubexpNames() {
		if i == 0 {
			continue
		}

		if name == "" {
			continue
		}

		playerData[name] = regexMatch[i]
	}

	return playerData, nil
}

func (h *HTTPHandler) UnparseMessageData(data map[string]string) string {
	var str strings.Builder

	for i, name := range h.MessageRegex.SubexpNames() {
		if i == 0 {
			continue
		}

		if name == "" {
			continue
		}

		if data[name] == "" {
			continue
		}

		str.WriteString(name)
		str.WriteString(data[name])
	}

	return str.String()
}
