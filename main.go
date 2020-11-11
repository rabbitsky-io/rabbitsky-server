package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"

	rsBot "rabbitsky/src/bot"
	rsChannel "rabbitsky/src/channel"
	rsHTTPHandler "rabbitsky/src/httphandler"
	rsWebSocket "rabbitsky/src/websocket"
)

func main() {
	/* Flag */
	serverPort := flag.Int("port", 8080, "Port the server will running on.")
	serverTick := flag.Int("tick", 10, "Server tick in Hz, more tick equal more bandwidth. Best to leave it default.")
	maxPlayers := flag.Int("max-players", 100, "Maximum Players in the server.")
	origin := flag.String("origin", "https://demo.rabbitsky.io", "Hostname of the website you will serve the Static HTML. Please do remember to insert http:// or https:// for this to works.")
	serverPassword := flag.String("admin-password", "", "Admin password to allow user to use admin command. Set on chat using '/admin [password]'. Leaving it empty will make command unusable.")
	addBots := flag.Int("add-bots", 0, "Spawn this number of bots to debug.")
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
	httpHandler, err := rsHTTPHandler.Init(channel, webSocket, *origin, *serverPassword)
	if err != nil {
		log.Fatalln("[Error] Could not initialize WebSocket. Err: ", err)
	}

	http.HandleFunc("/channel/players", httpHandler.GetChannelPlayers)
	http.HandleFunc("/channel/join", httpHandler.ChannelJoin)

	log.Println("Running on port", *serverPort)

	serverListen := fmt.Sprintf(":%d", *serverPort)
	log.Fatal(http.ListenAndServe(serverListen, nil))
}
