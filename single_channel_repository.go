package simplechatbot

import (
	models "github.com/c-rainbow/simplechatbot/models"
)

// SingleChannelRepository repo for a single channel
type SingleChannelRepository struct {
	chanInfo        *models.Channel
	commandsByBotID map[int64][]*models.Command
	baseRepo        *BaseRepository
}

// GetNewSingleChannelRepository returns a new single-channel repository
func GetNewSingleChannelRepository(chanInfo *models.Channel, baseRepo *BaseRepository) {
	commandsByBotID := map[int64][]*models.Command{}
	for _, command := range chanInfo.Commands {
		commandsByBotID[command.BotID] = append(commandsByBotID[command.BotID], &command)
	}
	repo := &SingleChannelRepository{
		chanInfo: chanInfo, commandsByBotID: commandsByBotID, baseRepo: baseRepo}
}

// GetAllBotIDs returns IDs of all bots associated with this channel
func (repo *SingleChannelRepository) GetAllBotIDs() []int64 {
	// Copy BotIDs to prevent unexpected modification
	copiedIDs := append([]int64{}, repo.chanInfo.BotIDs...)
	return copiedIDs
}

/*
func (repo *SingleChannelRepository) GetCommands []*models.Command {
	return nil
}
*/
// TODO:
// Refresh()  // refresh channel info, especially when command is updated
// AddCommand(command Command)  // all fields should be populated
// EditCommand(command Command)  // all fields should be populated
// DeleteCommand(command Command)  // command name is sufficient
