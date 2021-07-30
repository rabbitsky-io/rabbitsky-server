package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"

	rsBot "rabbitsky/src/bot"
	rsChannel "rabbitsky/src/channel"
	rsConfig "rabbitsky/src/config"
	rsHTTPHandler "rabbitsky/src/httphandler"
	rsWebSocket "rabbitsky/src/websocket"
)

func main() {

	/* Init Config */
	err := rsConfig.Init("config/config.json")
	if err != nil {
		log.Fatalln("[Error] Could not initialize WebSocket. Err: ", err)
	}
	//log.Println(rsConfig.DebugConfig())

	/* Flag */
	serverPort := flag.Int("port", rsConfig.Cfg.ServerPort, "Port the server will running on.")
	serverTick := flag.Int("tick", rsConfig.Cfg.ServerTick, "Server tick in Hz, more tick equal more bandwidth. Best to leave it default.")
	maxPlayers := flag.Int("max-players", rsConfig.Cfg.MaxPlayers, "Maximum Players in the server.")
	origin := flag.String("origin", rsConfig.Cfg.Origin, "Hostname of the website you will serve the Static HTML. Please do remember to insert http:// or https:// for this to works.")
	serverPassword := flag.String("admin-password", rsConfig.Cfg.ServerPassword, "Admin password to allow user to use admin command. Set on chat using '/admin [password]'. Leaving it empty will make command unusable.")
	limitPosMin := flag.String("limit-position-min", rsConfig.Cfg.LimitPosMin, "Position (X,Y,Z) Minimum that user allowed. If user position is less than this value, they will be disconnected.")
	limitPosMax := flag.String("limit-position-max", rsConfig.Cfg.LimitPosMax, "Position (X,Y,Z) Maximum that user allowed. If user position is more than this value, they will be disconnected.")
	addBots := flag.Int("add-bots", rsConfig.Cfg.DefaultBotCount, "Spawn this number of bots to debug.")
	flag.Parse()

	/* Init WebSocket */
	webSocket, err := rsWebSocket.Init(*origin)
	if err != nil {
		log.Fatalln("[Error] Could not initialize WebSocket. Err: ", err)
	}

	/* Init Channel */
	channel, err := rsChannel.Init(*maxPlayers, *serverTick)
	if err != nil {
		log.Fatalln("[Error] Could not initialize WebSocket. Err: ", err)
	}

	/* Init Bot If Not 0 */
	go rsBot.AddBot(channel, *addBots, *serverTick)

	/* Init HTTP Handler */
	httpHandler, err := rsHTTPHandler.Init(channel, webSocket, *origin, *serverPassword, *limitPosMin, *limitPosMax)
	if err != nil {
		log.Fatalln("[Error] Could not initialize WebSocket. Err: ", err)
	}

	http.HandleFunc("/channel/players", httpHandler.GetChannelPlayers)
	http.HandleFunc("/channel/join", httpHandler.ChannelJoin)

	log.Println("Running on port", *serverPort)

	serverListen := fmt.Sprintf(":%d", *serverPort)
	log.Fatal(http.ListenAndServe(serverListen, nil))
}
