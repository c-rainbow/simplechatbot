package commandplugins

import (
	"github.com/c-rainbow/simplechatbot/client"
	models "github.com/c-rainbow/simplechatbot/models"
	chatplugins "github.com/c-rainbow/simplechatbot/plugins/chat"
	"github.com/c-rainbow/simplechatbot/repository"
	twitch_irc "github.com/gempir/go-twitch-irc"
)

const (
	DeleteCommandPluginType = "DeleteCommandPluginType"
)

type DeleteCommandPlugin struct {
	ircClient client.TwitchClientT
	repo      repository.SingleBotRepositoryT
}

var _ chatplugins.ChatCommandPluginT = (*DeleteCommandPlugin)(nil)

func NewDeleteCommandPlugin(
	ircClient client.TwitchClientT, repo repository.SingleBotRepositoryT) chatplugins.ChatCommandPluginT {
	return &DeleteCommandPlugin{ircClient: ircClient, repo: repo}
}

func (plugin *DeleteCommandPlugin) GetPluginType() string {
	return DeleteCommandPluginType
}

func (plugin *DeleteCommandPlugin) Run(
	command *models.Command, channel string, sender *twitch_irc.User, message *twitch_irc.Message) error {
	// Validate inputs, action, and output
	commandToDelete, err := plugin.ValidateInputs(command, channel, DeleteCommandPluginType, sender, message)
	converted, err := plugin.Action(plugin.repo.GetBotInfo(), command, commandToDelete, channel, sender, message, err)
	return CommonOutput(plugin.ircClient, channel, converted, err)
}

func (plugin *DeleteCommandPlugin) ValidateInputs(
	command *models.Command, channel string, expectedPluginType string, sender *twitch_irc.User,
	message *twitch_irc.Message) (*models.Command, error) {

	err := CommonValidateInputs(command, channel, DeleteCommandPluginType, sender, message)
	if err != nil {
		return nil, err
	}

	// Parse chat message to command struct
	commandNameToDelete, _ := GetCommandNameAndResponseTextFromChat(message.Text)
	commandToDelete := plugin.repo.GetCommandByChannelAndName(channel, commandNameToDelete)
	if commandToDelete == nil {
		return nil, chatplugins.CommandNotFoundError
	}
	return commandToDelete, nil
}

func (plugin *DeleteCommandPlugin) Action(
	botInfo *models.Bot, command *models.Command, commandToDelete *models.Command, channel string,
	sender *twitch_irc.User, message *twitch_irc.Message, err error) (string, error) {
	// Check what error is returned from ValidateInput()
	if err != nil {
		return "Failed to validate input", err
	}

	err = plugin.repo.DeleteCommand(channel, commandToDelete)
	if err != nil {
		return "Failed to delete command", err
	}

	args := []string{commandToDelete.Name}
	converted, err := CommonConvertToResponseText(command, channel, sender, message, args, err)
	return converted, err
}

type DeleteCommandPluginFactory struct {
	ircClient client.TwitchClientT
	repo      repository.SingleBotRepositoryT
}

var _ chatplugins.ChatCommandPluginFactoryT = (*DeleteCommandPluginFactory)(nil)

func NewDeleteCommandPluginFactory(
	ircClient client.TwitchClientT, repo repository.SingleBotRepositoryT) chatplugins.ChatCommandPluginFactoryT {
	return &DeleteCommandPluginFactory{ircClient: ircClient, repo: repo}
}

func (plugin *DeleteCommandPluginFactory) GetPluginType() string {
	return DeleteCommandPluginType
}

func (plugin *DeleteCommandPluginFactory) BuildNewPlugin() chatplugins.ChatCommandPluginT {
	return NewDeleteCommandPlugin(plugin.ircClient, plugin.repo)
}

func (plugin *DeleteCommandPlugin) GetResponseKey(err error) string {
	// Normal case.
	if err == nil {
		return models.DefaultResponseKey
	}

	/*
		Failure case.
		Design decision: We can return different messages per error type in two ways
		(1) switch statement with each known error cases, manually assigning response key, like the code below
		(2) Each error has unique error type string, and we use it as response key.
			For example, parsedResponse, exists := command.Responsees[NoPermissionError.Key()]

		Design 1 was chosen because mapping between error type and error message is not necessarily 1-to-1.
		Multiple error types can produce the same error message. With design 1, it's also more obvious in the code
		to see which error is connected to which message key.
	*/
	switch err {
	case chatplugins.ErrCommandNotFound: // Command name is not found. Likely syncronization issue
		fallthrough
	case chatplugins.ErrNoPermission: // User has no permission
		fallthrough
	case chatplugins.ErrNotEnoughArguments: // Arguments
		fallthrough
	case chatplugins.ErrTargetCommandAlreadyExists: // Target command already exists and cannot be added
		fallthrough
	case chatplugins.ErrTargetCommandNotFound: // This is not relevant to AddCommand
		fallthrough
	default:
		return models.DefaultFailureResponseKey
	}
}
