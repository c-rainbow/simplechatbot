package commandplugins

import (
	models "github.com/c-rainbow/simplechatbot/models"
	chatplugins "github.com/c-rainbow/simplechatbot/plugins/chat"
	"github.com/c-rainbow/simplechatbot/repository"
)

// Returns command model if the target command name exists.
func getTargetCommand(
	channel string, targetName string, repo repository.SingleBotRepositoryT) (*models.Command, error) {

	if targetName == "" {
		return nil, chatplugins.ErrNotEnoughArguments
	}
	targetCommand := repo.GetCommandByChannelAndName(channel, targetName)
	return targetCommand, nil
}
