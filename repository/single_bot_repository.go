package repository

import (
	"fmt"

	models "github.com/c-rainbow/simplechatbot/models"
)

// SingleBotRepositoryT interface for single bot repository
type SingleBotRepositoryT interface {
	GetBotInfo() *models.Bot
	GetAllChannels() []*models.Channel
	GetCommandByChannelAndName(channel string, commandName string) *models.Command
	AddCommand(channel string, command *models.Command) error
	EditCommand(channel string, command *models.Command) error
	DeleteCommand(channel string, command *models.Command) error
	ListCommands(channel string) ([]*models.Command, error)
}

// SingleBotRepository repository for a single bot
type SingleBotRepository struct {
	botInfo  *models.Bot
	baseRepo BaseRepositoryT
}

var _ SingleBotRepositoryT = (*SingleBotRepository)(nil)

// NewSingleBotRepository creates a new single bor repository with given bot and base repo.
func NewSingleBotRepository(botInfo *models.Bot, baseRepo BaseRepositoryT) *SingleBotRepository {
	return &SingleBotRepository{botInfo: botInfo, baseRepo: baseRepo}
}

// GetBotInfo returns bot model struct
func (repo *SingleBotRepository) GetBotInfo() *models.Bot {
	return repo.botInfo
}

// GetAllChannels returns all channels for this bot.
func (repo *SingleBotRepository) GetAllChannels() []*models.Channel {
	botID := repo.botInfo.TwitchID
	channels := repo.baseRepo.GetAllChannelsForBot(botID)
	return channels
}

// GetCommandByChannelAndName gets commands by channel name and command name
func (repo *SingleBotRepository) GetCommandByChannelAndName(
	channel string, commandName string) *models.Command {
	return repo.baseRepo.GetCommand(repo.botInfo.TwitchID, channel, commandName)
}

// AddCommand adds command.
func (repo *SingleBotRepository) AddCommand(channel string, command *models.Command) error {
	err := repo.baseRepo.AddCommand(channel, command)
	if err != nil {
		fmt.Println("Error from baseRepo, AddCommand: ", err.Error())
		return err
	}
	return nil
}

// EditCommand edits command
func (repo *SingleBotRepository) EditCommand(channel string, command *models.Command) error {
	err := repo.baseRepo.EditCommand(channel, command)
	if err != nil {
		fmt.Println("Error from baseRepo, EditCommand: ", err.Error())
		return err
	}
	return nil
}

// DeleteCommand deletes command
func (repo *SingleBotRepository) DeleteCommand(channel string, command *models.Command) error {
	err := repo.baseRepo.DeleteCommand(channel, command)
	if err != nil {
		fmt.Println("Error from baseRepo, DeleteCommand: ", err.Error())
		return err
	}
	return nil
}

// ListCommands lists command
// TODO: Eventually, point this to a web page
func (repo *SingleBotRepository) ListCommands(channel string) ([]*models.Command, error) {
	commands, err := repo.baseRepo.ListCommands(channel)
	if err != nil {
		fmt.Println("Error from baseRepo, ListCommands: ", err.Error())
		return nil, err
	}
	return commands, nil
}
