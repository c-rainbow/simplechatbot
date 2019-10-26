package commandplugins

import (
	"errors"
	"log"

	"github.com/c-rainbow/simplechatbot/client"
	models "github.com/c-rainbow/simplechatbot/models"
	"github.com/c-rainbow/simplechatbot/parser"
	chatplugins "github.com/c-rainbow/simplechatbot/plugins/chat"
	"github.com/c-rainbow/simplechatbot/plugins/chat/common"
	"github.com/c-rainbow/simplechatbot/repository"
	twitch_irc "github.com/gempir/go-twitch-irc"
)

const (
	EditCommandPluginType = "EditCommandPluginType"
)

type EditCommandPluginFactory struct {
	ircClient client.TwitchClientT
	repo      repository.SingleBotRepositoryT
}

var _ chatplugins.ChatCommandPluginFactoryT = (*EditCommandPluginFactory)(nil)

func NewEditCommandPluginFactory(
	ircClient client.TwitchClientT, repo repository.SingleBotRepositoryT) chatplugins.ChatCommandPluginFactoryT {
	return &EditCommandPluginFactory{ircClient: ircClient, repo: repo}
}

func (plugin *EditCommandPluginFactory) GetPluginType() string {
	return EditCommandPluginType
}

func (plugin *EditCommandPluginFactory) BuildNewPlugin() chatplugins.ChatCommandPluginT {
	return NewEditCommandPlugin(plugin.ircClient, plugin.repo)
}

type EditCommandPlugin struct {
	ircClient client.TwitchClientT
	repo      repository.SingleBotRepositoryT
}

var _ chatplugins.ChatCommandPluginT = (*EditCommandPlugin)(nil)

func NewEditCommandPlugin(
	ircClient client.TwitchClientT, repo repository.SingleBotRepositoryT) chatplugins.ChatCommandPluginT {
	return &EditCommandPlugin{ircClient: ircClient, repo: repo}
}

func (plugin *EditCommandPlugin) GetPluginType() string {
	return EditCommandPluginType
}

func (plugin *EditCommandPlugin) ReactToChat(
	command *models.Command, channel string, sender *twitch_irc.User, message *twitch_irc.PrivateMessage) {
	var err error
	var targetCommand *models.Command

	targetCommandName, targetResponse := GetTargetCommandNameAndResponse(message.Message)

	// TODO: Is it possible to get away from this continuous err == nil check?
	err = common.ValidateBasicInputs(command, channel, EditCommandPluginType, sender, message)
	if err == nil {
		targetCommand, err = plugin.GetTargetCommand(channel, targetCommandName)
	}
	if err == nil {
		err = plugin.EditTargetCommand(targetCommand, targetResponse)
	}
	if err == nil {
		err = plugin.ValidateTargetCommand(targetCommand, targetResponse)
	}
	if err == nil {
		err = plugin.repo.EditCommand(channel, targetCommand)
	}

	responseText, err := plugin.GetResponseText(command, targetCommand, channel, sender, message, err)

	common.SendToChatClient(plugin.ircClient, channel, responseText)
	common.HandleError(err)
}

func (plugin *EditCommandPlugin) GetTargetCommand(channel string, targetName string) (*models.Command, error) {
	if targetName == "" {
		return nil, chatplugins.ErrNotEnoughArguments
	}
	targetCommand := plugin.repo.GetCommandByChannelAndName(channel, targetName)
	return targetCommand, nil
}

// This function is slightly different between add/edit/delete command. Hard to merge into a common function.
func (plugin *EditCommandPlugin) ValidateTargetCommand(targetCommand *models.Command, targetResponse string) error {
	// Can't edit non-existing command
	if targetCommand == nil {
		return chatplugins.ErrTargetCommandNotFound
	}
	if targetResponse == "" {
		return chatplugins.ErrNotEnoughArguments
	}
	if targetCommand.Responses[models.DefaultResponseKey].RawText != targetResponse {
		return errors.New("Default response different from the input")
	}
	return nil
}

// Only needed for EditCommand. Update the default response of the command model
func (plugin *EditCommandPlugin) EditTargetCommand(targetCommand *models.Command, targetResponse string) error {

	// Build default response with the given message.
	defaultResponse := parser.ParseResponse(targetResponse)
	err := parser.Validate(defaultResponse)
	if err != nil {
		log.Println("Failed to validate target response")
		return err
	}

	targetCommand.Responses[models.DefaultResponseKey] = *defaultResponse
	return nil
}

// Get response text of the executed command, based on the errors and progress so far.
func (plugin *EditCommandPlugin) GetResponseText(
	command *models.Command, targetCommand *models.Command, channel string, sender *twitch_irc.User,
	message *twitch_irc.PrivateMessage, err error) (string, error) {
	// Get response key, build args, get parsed response, and convert it to text
	responseKey := plugin.GetResponseKey(err)
	args := []string{targetCommand.Name}
	return common.ConvertToResponseText(command, responseKey, channel, sender, message, args)
}

func (plugin *EditCommandPlugin) GetResponseKey(err error) string {
	// Normal case.
	if err == nil {
		return models.DefaultResponseKey
	}
	// Failure cases
	switch err {
	case chatplugins.ErrCommandNotFound: // Command name is not found. Likely synchronization issue
		fallthrough
	case chatplugins.ErrNoPermission: // User has no permission
		fallthrough
	case chatplugins.ErrNotEnoughArguments: // Arguments
		fallthrough
	case chatplugins.ErrTargetCommandNotFound: // Target command not found and cannot be edited
		fallthrough
	default:
		return models.DefaultFailureResponseKey
	}
}
