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
	SelfBanPluginType = "SelfBanPluginType"
	DefaultBanSeconds = 3
)

type SelfBanPluginFactory struct {
	ircClient client.TwitchClientT
}

var _ chatplugins.ChatCommandPluginFactoryT = (*SelfBanPluginFactory)(nil)

func NewSelfBanPluginFactory(ircClient client.TwitchClientT) chatplugins.ChatCommandPluginFactoryT {
	return &SelfBanPluginFactory{ircClient: ircClient}
}

func (plugin *SelfBanPluginFactory) GetPluginType() string {
	return SelfBanPluginType
}

func (plugin *SelfBanPluginFactory) BuildNewPlugin() chatplugins.ChatCommandPluginT {
	return NewSelfBanPlugin(plugin.ircClient)
}

type SelfBanPlugin struct {
	ircClient client.TwitchClientT
	repo      repository.SingleBotRepositoryT
}

var _ chatplugins.ChatCommandPluginT = (*SelfBanPlugin)(nil)

func NewSelfBanPlugin(ircClient client.TwitchClientT) chatplugins.ChatCommandPluginT {
	return &SelfBanPlugin{ircClient: ircClient}
}

func (plugin *SelfBanPlugin) GetPluginType() string {
	return SelfBanPluginType
}

func (plugin *SelfBanPlugin) ReactToChat(
	command *models.Command, channel string, sender *twitch_irc.User, message *twitch_irc.Message) {
	var banTime int
	// TODO: Is it possible to get away from this continuous err == nil check?
	err := common.ValidateBasicInputs(command, channel, SelfBanPluginType, sender, message)
	if err == nil {
		//banTime, err := plugin.ParseTime(message.Text)
		banTime = DefaultBanSeconds
	}
	if err != nil {
		banTime = DefaultBanSeconds
	}

	responseText, err := plugin.GetResponseText(command, channel, banTime, sender, message, err)

	plugin.TryBanUser(channel, sender, banTime, responseText)
	common.HandleError(err)
}

func (plugin *SelfBanPlugin) TryBanUser(channel string, sender *twitch_irc.User, banTime int, responseText string) {
	plugin.ircClient.Say(channel, "/timeout "+sender.Username+" "+strconv.Itoa(banTime))
	plugin.ircClient.Say(channel, responseText)
}

// Get response text of the executed command, based on the errors and progress so far.
func (plugin *SelfBanPlugin) GetResponseText(
	command *models.Command, channel string, banTime int, sender *twitch_irc.User, message *twitch_irc.Message,
	err error) (string, error) {

	// Get response key, build args, get parsed response, and convert it to text
	responseKey := plugin.GetResponseKey(err)
	return common.ConvertToResponseText(command, responseKey, channel, sender, message, nil)
}

func (plugin *SelfBanPlugin) GetResponseKey(err error) string {
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
