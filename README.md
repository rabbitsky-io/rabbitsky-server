# Rabbit Sky - Server Side

[![Rabbit Sky Server Status](https://circleci.com/gh/rabbitsky-io/rabbitsky-server.svg?style=shield)](<LINK>)

This is the Server Side / Part Two of the Rabbit Sky. It job is to run WebSocket server to connects all players.

We recommend using Linux for the server side. Works well on Linux (Tested on CentOS 7 and Ubuntu 18), also tested on Windows, but haven't tested it on Mac yet since we don't have one.

[Click here to see the Static Web / Part One](https://github.com/rabbitsky-io/rabbitsky-web)

## Download
Please refer to [Release Page](https://github.com/rabbitsky-io/rabbitsky-server/releases) to download the latest binary.

## Running The Server
`./rabbitsky [-port num] [-tick num] [-max-players num] [-origin string]`

## Command Parameter
| Args     | Type | Default | Description |
| -------- | ---- | ------- | ----------- |
| -port | number | 8080 | Port this app will listen to. |
| -tick | number | 10 | Server tick in Hz. How often the client and server communicate to each other. Please note that increasing the value will increase the usage of bandwidth for both client and server. We recommend to not increase or decrease the value. Increasing it to more than 30 can cause some problems, and we do not allow it to be set more than 60. |
| -max-players | number | 100 | Maximum players this server can serve. Please note that the more players in the server, the more bandwidth is used by both the client and the server. Default is the best for both, increasing it to more than 250 can cause some problems.|
| -origin | string | https://demo.rabbitsky.io | URL of the domain you use for the static file (rabbitsky-web). This is used for CORS. Please remember to input the scheme (http:// or https://) and remove trailing slash. |

## Debug Parameter
| Args     | Type | Default | Description |
| -------- | ---- | ------- | ----------- |
| -add-bots | number | 0 | Add x numbers of bots in the server. Bot will randomly moving. It's for debugging only and it is recommended to not use this in production. |

## Always Online
Sometime the server can crash during the event, so if you want it to be automatically restart when crash, you can create a script or service.

If you're using Linux we recommend you to create your own service using systemd. Please refer to this post to start: [How to create systemd service?](https://linuxconfig.org/how-to-create-systemd-service-unit-in-linux).

## Making Your Own Build
It's relatively easy. We created a Makefile just for that! All you have to do is to have golang on your PC / Mac. [Click here to see how to install golang](https://golang.org/doc/install).

If you have golang already, you can try command `make` on Linux / Mac. If you're using another OS, you can use command `go build .`

If it works, you can test your build with command `./rabbitsky -help`

# Security
We haven't found a perfect way to secure the WebSocket connection, except using Origin header.

Token is the only thing we consider during the creation of the server, but it only best with user authentication / login. Generating token for each non-authenticated user is kinda useless, because everyone can create a token easily, so it's easy to break in. Also we don't really want to use Cookie.

Let me know if you have an idea for this matter!

## Donate
Liking this and having a spare of money? Consider donating!

[![Donate](https://www.paypalobjects.com/en_US/i/btn/btn_donateCC_LG.gif)](https://paypal.me/wibisaja)

## Contributing

Yes

## 3rd Party Module

[catinello/base62](github.com/catinello/base62)
[gorilla/websocket](github.com/gorilla/websocket)
[orcaman/concurrent-map](github.com/orcaman/concurrent-map)