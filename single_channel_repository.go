package simplechatbot

import (
	models "github.com/c-rainbow/simplechatbot/models"
)

// SingleChannelRepository repo for a single channel
type SingleChannelRepository struct {
	chanInfo *models.Channel
	baseRepo *BaseRepository
}

func (repo *SingleChannelRepository) GetAllBots() []int64 {
	return repo.chanInfo.BotIDs
}

func (repo *SingleChannelRepository) GetCommandsByBotID() {

}

func (repo *SingleChannelRepository) UpdateCommand() {

}
