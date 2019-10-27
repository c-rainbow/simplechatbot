package botserver

import (
	"github.com/c-rainbow/simplechatbot/botmanager"
	"github.com/c-rainbow/simplechatbot/restserver"
)

// Bot Server. It has BotManager and REST server.
// One Bot Server should have onyl one base repository, which is in turn used by the bot manager.

type BotServerT interface {
	Start()
	Shutdown()
}

type BotServer struct {
	restServer restserver.RestServerT
	botManager botmanager.BotManagerT
}

var _ BotServerT = (*BotServer)(nil)

func NewBotServer(restServer restserver.RestServerT, botManager botmanager.BotManagerT) BotServerT {
	return &BotServer{restServer: restServer, botManager: botManager}
}

func (server *BotServer) Start() {
	if server.restServer != nil {
		server.restServer.Start()
	}
	server.botManager.Start()
}

func (server *BotServer) Shutdown() {
	if server.restServer != nil {
		server.restServer.Shutdown()
	}
	server.botManager.Shutdown()
}
