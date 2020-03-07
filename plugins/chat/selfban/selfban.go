package selfban

import (
	"strconv"

	"github.com/c-rainbow/simplechatbot/client"
	models "github.com/c-rainbow/simplechatbot/models"
	chatplugins "github.com/c-rainbow/simplechatbot/plugins/chat"
	"github.com/c-rainbow/simplechatbot/plugins/chat/common"
	"github.com/c-rainbow/simplechatbot/repository"
	twitch_irc "github.com/gempir/go-twitch-irc"
)

const (
	// BanOneselfPluginType plugin type name to enable chatters to ban themselves
	BanOneselfPluginType = "BanOneselfPluginType"
	defaultBanSeconds    = 3
)

// BanOneselfPlugin plugin to enable chatters to ban themselves
type BanOneselfPlugin struct {
	ircClient client.TwitchClientT
	repo      repository.SingleBotRepositoryT
}

var _ chatplugins.ChatCommandPluginT = (*BanOneselfPlugin)(nil)

// NewBanOneselfPlugin creates a new BanOneselfPlugin
func NewBanOneselfPlugin(ircClient client.TwitchClientT) chatplugins.ChatCommandPluginT {
	return &BanOneselfPlugin{ircClient: ircClient}
}

// GetPluginType get splugin type
func (plugin *BanOneselfPlugin) GetPluginType() string {
	return BanOneselfPluginType
}

// ReactToChat reacts to chat
func (plugin *BanOneselfPlugin) ReactToChat(
	command *models.Command, channel string, sender *twitch_irc.User, message *twitch_irc.PrivateMessage) {
	var banTime int
	// TODO: Is it possible to get away from this continuous err == nil check?
	err := common.ValidateBasicInputs(command, channel, BanOneselfPluginType, sender, message)
	if err == nil {
		//banTime, err := plugin.ParseTime(message.Message)
		banTime = defaultBanSeconds
	}
	if err != nil {
		banTime = defaultBanSeconds
	}

	responseText, err := plugin.getResponseText(command, channel, banTime, sender, message, err)

	plugin.tryBanUser(channel, sender, banTime, responseText)
	common.HandleError(err)
}

func (plugin *BanOneselfPlugin) tryBanUser(channel string, sender *twitch_irc.User, banTime int, responseText string) {
	plugin.ircClient.Say(channel, "/timeout "+sender.Name+" "+strconv.Itoa(banTime))
}

// Get response text of the executed command, based on the errors and progress so far.
func (plugin *BanOneselfPlugin) getResponseText(
	command *models.Command, channel string, banTime int, sender *twitch_irc.User, message *twitch_irc.PrivateMessage,
	err error) (string, error) {

	// Get response key, build args, get parsed response, and convert it to text
	responseKey := plugin.getResponseKey(err)
	return common.ConvertToResponseText(command, responseKey, channel, sender, message, nil)
}

func (plugin *BanOneselfPlugin) getResponseKey(err error) string {
	// Normal case.
	if err == nil {
		return models.DefaultResponseKey
	}
	// Failure case.
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
