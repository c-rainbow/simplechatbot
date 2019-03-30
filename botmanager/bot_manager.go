package botmanager

import (
	"github.com/c-rainbow/simplechatbot/bot"
)

type BotManagerT interface {
	Start()
	Shutdown()
}

type BotManager struct {
	bots []bot.TwitchChatBot
}

var _ BotManagerT = (*BotManager)(nil)

func (manager *BotManager) Start() {
	for _, bot := range manager.bots {
		bot.Start()
	}
}

func (manager *BotManager) Shutdown() {
	for _, bot := range manager.bots {
		bot.Shutdown()
	}
}
