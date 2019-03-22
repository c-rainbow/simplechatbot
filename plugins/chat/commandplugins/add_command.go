package commandplugins

import (
	"log"
	"strconv"

	"github.com/c-rainbow/simplechatbot/client"
	models "github.com/c-rainbow/simplechatbot/models"
	parser "github.com/c-rainbow/simplechatbot/parser"
	chatplugins "github.com/c-rainbow/simplechatbot/plugins/chat"
	"github.com/c-rainbow/simplechatbot/plugins/chat/common"
	"github.com/c-rainbow/simplechatbot/repository"
	twitch_irc "github.com/gempir/go-twitch-irc"
)

const (
	AddCommandPluginType = "AddCommandPluginType"
)

type AddCommandPluginFactoryT struct {
	ircClient client.TwitchClientT
	repo      repository.SingleBotRepositoryT
}

var _ chatplugins.ChatCommandPluginFactoryT = (*AddCommandPluginFactoryT)(nil)

func NewAddCommandPluginFactory(
	ircClient client.TwitchClientT, repo repository.SingleBotRepositoryT) chatplugins.ChatCommandPluginFactoryT {
	return &AddCommandPluginFactoryT{ircClient: ircClient, repo: repo}
}

func (plugin *AddCommandPluginFactoryT) GetPluginType() string {
	return AddCommandPluginType
}

func (plugin *AddCommandPluginFactoryT) BuildNewPlugin() chatplugins.ChatCommandPluginT {
	return NewAddCommandPlugin(plugin.ircClient, plugin.repo)
}

type AddCommandPlugin struct {
	ircClient client.TwitchClientT
	repo      repository.SingleBotRepositoryT
}

var _ chatplugins.ChatCommandPluginT = (*AddCommandPlugin)(nil)

func NewAddCommandPlugin(
	ircClient client.TwitchClientT, repo repository.SingleBotRepositoryT) chatplugins.ChatCommandPluginT {
	return &AddCommandPlugin{ircClient: ircClient, repo: repo}
}

func (plugin *AddCommandPlugin) GetPluginType() string {
	return AddCommandPluginType
}

func (plugin *AddCommandPlugin) ReactToChat(
	command *models.Command, channel string, sender *twitch_irc.User, message *twitch_irc.Message) {
	var err error
	var targetCommand *models.Command

	targetCommandName, targetResponse := GetTargetCommandNameAndResponse(message.Text)

	// TODO: Is it possible to get away from this continuous err == nil check?
	err = common.ValidateBasicInputs(command, channel, AddCommandPluginType, sender, message)
	if err == nil {
		targetCommand, err = plugin.GetTargetCommand(channel, targetCommandName)
	}
	if err == nil {
		err = plugin.ValidateTargetCommand(targetCommand, targetResponse)
	}
	if err == nil {
		targetCommand, err = plugin.BuildTargetCommand(targetCommandName, targetResponse, message.ChannelID)
	}
	if err == nil {
		err = plugin.repo.AddCommand(channel, targetCommand)
	}

	responseText, err := plugin.GetResponseText(command, targetCommand, channel, sender, message, err)

	common.SendToChatClient(plugin.ircClient, channel, responseText)
	common.HandleError(err)
}

func (plugin *AddCommandPlugin) GetTargetCommand(channel string, targetName string) (*models.Command, error) {
	if targetName == "" {
		return nil, chatplugins.ErrNotEnoughArguments
	}
	targetCommand := plugin.repo.GetCommandByChannelAndName(channel, targetName)
	return targetCommand, nil
}

// This function is slightly different between add/edit/delete command. Hard to merge into a common function.
func (plugin *AddCommandPlugin) ValidateTargetCommand(targetCommand *models.Command, targetResponse string) error {
	// Can't add already existing command
	if targetCommand != nil {
		return chatplugins.ErrTargetCommandAlreadyExists
	}
	if targetResponse == "" {
		return chatplugins.ErrNotEnoughArguments
	}
	return nil
}

// Only needed for AddCommand
func (plugin *AddCommandPlugin) BuildTargetCommand(
	targetCommandName string, targetResponse string, channelIDStr string) (*models.Command, error) {
	// Convert channelIDstr to int
	channelID, err := strconv.Atoi(channelIDStr)
	if err != nil {
		// This error statement means ChannelID != channel's TwitchID, or a bug with IRC library
		log.Println("Failed to convert ChannelID '", channelIDStr, "' to int.")
		return nil, chatplugins.ErrInvalidArgument
	}

	// Build default response with the given message.
	defaultResponse := parser.ParseResponse(targetResponse)
	err = parser.Validate(defaultResponse)
	if err != nil {
		log.Println("Failed to validate target response")
		return nil, err
	}

	// TODO: Find a nice, descriptive failure message.
	failureRespopnse := parser.ParseResponse(chatplugins.DefaultFailureMessage)
	err = parser.Validate(failureRespopnse)
	if err != nil {
		log.Println("Failed to validate default failure response")
		return nil, err
	}

	// Build a new Command object. Other fields are auto-initialized.
	botID := plugin.repo.GetBotInfo().TwitchID
	targetCommand := models.NewCommand(
		botID, int64(channelID), targetCommandName, CommandResponsePluginType, defaultResponse, failureRespopnse)
	return targetCommand, nil
}

// Get response text of the executed command, based on the errors and progress so far.
func (plugin *AddCommandPlugin) GetResponseText(
	command *models.Command, targetCommand *models.Command, channel string, sender *twitch_irc.User,
	message *twitch_irc.Message, err error) (string, error) {
	// Get response key, build args, get parsed response, and convert it to text
	responseKey := plugin.GetResponseKey(err)
	args := []string{""}
	if targetCommand != nil {
		args = []string{targetCommand.Name}
	}
	return common.ConvertToResponseText(command, responseKey, channel, sender, message, args)
}

func (plugin *AddCommandPlugin) GetResponseKey(err error) string {
	// Normal case.
	if err == nil {
		return models.DefaultResponseKey
	}

	/*
		Failure case.
		Design decision: We can return different messages per error type in two ways
		(1) switch statement with each known error cases, manually assigning response key, like the code below
		(2) Each error has unique error type string, and we use it as response key. For example,
			parsedResponse, exists := Responses[NoPermissionError.Key()]

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
	default:
		return models.DefaultFailureResponseKey
	}
}
