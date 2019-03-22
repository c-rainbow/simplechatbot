package commandplugins

import (
	"github.com/c-rainbow/simplechatbot/client"
	models "github.com/c-rainbow/simplechatbot/models"
	chatplugins "github.com/c-rainbow/simplechatbot/plugins/chat"
	"github.com/c-rainbow/simplechatbot/repository"
	twitch_irc "github.com/gempir/go-twitch-irc"
)

const (
	EditCommandPluginType = "EditCommandPluginType"
)

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

func (plugin *EditCommandPlugin) Run(
	command *models.Command, channel string, sender *twitch_irc.User, message *twitch_irc.Message) error {
	repo := plugin.repo
	return CommonRun(repo, plugin.ircClient, EditCommandPluginType, repo.EditCommand, command,
		channel, sender, message)
}

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

func (plugin *EditCommandPlugin) GetResponseKey(err error) string {
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
