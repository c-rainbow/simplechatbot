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
	// AddCommandPluginType plugin type name to add new chat command of CommandResponsePluginType
	AddCommandPluginType = "AddCommandPluginType"
)

// AddCommandPlugin plugin to add new chat command of CommandResponsePluginType
// TODO: This plugin is called from chat message, and by default, the added command can be called by everyone.
type AddCommandPlugin struct {
	ircClient client.TwitchClientT
	repo      repository.SingleBotRepositoryT
}

var _ chatplugins.ChatCommandPluginT = (*AddCommandPlugin)(nil)

// NewAddCommandPlugin creates a new plugin
func NewAddCommandPlugin(
	ircClient client.TwitchClientT, repo repository.SingleBotRepositoryT) chatplugins.ChatCommandPluginT {
	return &AddCommandPlugin{ircClient: ircClient, repo: repo}
}

// GetPluginType returns plugin type
func (plugin *AddCommandPlugin) GetPluginType() string {
	return AddCommandPluginType
}

// ReactToChat reacts to chat
func (plugin *AddCommandPlugin) ReactToChat(
	command *models.Command, channel string, sender *twitch_irc.User, message *twitch_irc.PrivateMessage) {
	var err error
	var targetCommand *models.Command

	targetCommandName, targetResponse := GetTargetCommandNameAndResponse(message.Message)

	// TODO: Is it possible to get away from this continuous err == nil check?
	err = common.ValidateBasicInputs(command, channel, AddCommandPluginType, sender, message)
	if err == nil {
		targetCommand, err = getTargetCommand(channel, targetCommandName, plugin.repo)
	}
	if err == nil {
		err = plugin.validateTargetCommand(targetCommand, targetResponse)
	}
	if err == nil {
		// message.RoomID is integer Twitch channel ID. It used to be possible to have multiple rooms in a channel
		// (therefore multiple roomIDs), but this feature was gone on October 30, 2019.
		targetCommand, err = plugin.BuildTargetCommand(targetCommandName, targetResponse, message.RoomID)
	}
	if err == nil {
		err = plugin.repo.AddCommand(channel, targetCommand)
	}

	responseText, err := plugin.GetResponseText(command, targetCommand, channel, sender, message, err)

	common.SendToChatClient(plugin.ircClient, channel, responseText)
	common.HandleError(err)
}

// This function is slightly different between add/edit/delete command. Hard to merge into a common function.
func (plugin *AddCommandPlugin) validateTargetCommand(targetCommand *models.Command, targetResponse string) error {
	// Can't add already existing command
	if targetCommand != nil {
		return chatplugins.ErrTargetCommandAlreadyExists
	}
	if targetResponse == "" {
		return chatplugins.ErrNotEnoughArguments
	}
	return nil
}

// BuildTargetCommand builds command model from name. Only needed for AddCommand
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
	// TODO: dynamically get failure message from context (channel language, etc)
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

// GetResponseText gets response text of the executed command, based on the errors and progress so far.
func (plugin *AddCommandPlugin) GetResponseText(
	command *models.Command, targetCommand *models.Command, channel string, sender *twitch_irc.User,
	message *twitch_irc.PrivateMessage, err error) (string, error) {
	// Get response key, build args, get parsed response, and convert it to text
	responseKey := plugin.GetResponseKey(err)
	args := []string{""}
	if targetCommand != nil {
		args = []string{targetCommand.Name}
	}
	return common.ConvertToResponseText(command, responseKey, channel, sender, message, args)
}

// GetResponseKey returns response key from error type to build response text accordingly.
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
