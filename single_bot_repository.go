package simplechatbot

import (
	models "github.com/c-rainbow/simplechatbot/models"
)

type SingleBotRepositoryT interface {
	GetAllChannels() []*models.Channel
	GetCommandByChannelAndName(channel string, commandName string) *models.Command
}

// Repository for a single bot
type SingleBotRepository struct {
	botInfo  *models.Bot
	baseRepo BaseRepositoryT
}

var _ SingleBotRepositoryT = (*SingleBotRepository)(nil)

func NewSingleBotRepository(botInfo *models.Bot, baseRepo BaseRepositoryT) *SingleBotRepository {
	return &SingleBotRepository{botInfo: botInfo, baseRepo: baseRepo}
}

// GetAllChannels returns all channels for this bot.
func (repo *SingleBotRepository) GetAllChannels() []*models.Channel {
	botID := repo.botInfo.TwitchID
	channels := repo.baseRepo.GetAllChannelsForBot(botID)
	return channels
}

func (repo *SingleBotRepository) GetCommandByChannelAndName(
	channel string, commandName string) *models.Command {
	return repo.baseRepo.GetCommand(repo.botInfo.TwitchID, channel, commandName)
}
