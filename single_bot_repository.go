package simplechatbot

import (
	"fmt"

	models "github.com/c-rainbow/simplechatbot/models"
)

type SingleBotRepositoryT interface {
	GetAllChannels() []*models.Channel
	GetCommandByChannelAndName(channel string, commandName string) *models.Command
	AddCommand(channel string, command *models.Command) error
	EditCommand(channel string, command *models.Command) error
	DeleteCommand(channel string, command *models.Command) error
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

func (repo *SingleBotRepository) AddCommand(channel string, command *models.Command) error {
	err := repo.baseRepo.AddCommand(channel, command)
	if err != nil {
		fmt.Println("Error from baseRepo: ", err.Error())
		return err
	}
	return nil
}

func (repo *SingleBotRepository) EditCommand(channel string, command *models.Command) error {
	err := repo.baseRepo.EditCommand(channel, command)
	if err != nil {
		fmt.Println("Error from baseRepo: ", err.Error())
		return err
	}
	return nil
}

func (repo *SingleBotRepository) DeleteCommand(channel string, command *models.Command) error {
	err := repo.baseRepo.DeleteCommand(channel, command)
	if err != nil {
		fmt.Println("Error from baseRepo: ", err.Error())
		return err
	}
	return nil
}
