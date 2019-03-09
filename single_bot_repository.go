package simplechatbot

import (
	models "github.com/c-rainbow/simplechatbot/models"
	dynamo "github.com/guregu/dynamo"
)

// Repository for a single bot
type SingleBotRepository struct {
	botInfo *models.Bot
	db      *dynamo.DB
}

func (repo *SingleBotRepository) GetAllChannels() {

}
