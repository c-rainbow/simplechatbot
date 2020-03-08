package botmanager

// One bot manager has bots with the same base repository.

import (
	"github.com/c-rainbow/simplechatbot/bot"
	"github.com/c-rainbow/simplechatbot/repository"
)

type BotManagerT interface {
	Start()
	Shutdown()
	// AddBot(bot bot.TwitchChatBotT)?
	// RemoveBot?
}

type BotManager struct {
	repo repository.BaseRepositoryT // Is this really needed?
	bots []bot.TwitchChatBotT
}

var _ BotManagerT = (*BotManager)(nil)

func NewBotManagerFromRepository(baseRepo repository.BaseRepositoryT) BotManagerT {
	// TODO: Currently, there is 1-1 relationship between repository and bot manager.
	// In the future, one repository should be able to be used by multiple bot managers.
	botModels := baseRepo.GetAllBots()
	bots := make([]bot.TwitchChatBotT, len(botModels))
	for i, botModel := range botModels {
		bots[i] = bot.DefaultTwitchChatBot(botModel, baseRepo)
	}

	return &BotManager{repo: baseRepo, bots: bots}
}

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
