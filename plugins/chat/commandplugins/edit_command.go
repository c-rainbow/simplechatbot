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
