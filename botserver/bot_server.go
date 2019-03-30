package botserver

import (
	"github.com/c-rainbow/simplechatbot/botmanager"
	"github.com/c-rainbow/simplechatbot/restserver"
)

// Bot Server. It has BotManager and REST server

type BotServer struct {
	restServer restserver.RestServerT
	botManager botmanager.BotManagerT
}

func NewBotServer() *BotServer {

	return nil
}
