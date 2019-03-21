package commandplugins

import (
	"github.com/c-rainbow/simplechatbot/client"
	models "github.com/c-rainbow/simplechatbot/models"
	chatplugins "github.com/c-rainbow/simplechatbot/plugins/chat"
	"github.com/c-rainbow/simplechatbot/repository"
	twitch_irc "github.com/gempir/go-twitch-irc"
)

const (
	AddCommandPluginType = "AddCommandPluginType"
)

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

func (plugin *AddCommandPlugin) Run(
	command *models.Command, channel string, sender *twitch_irc.User, message *twitch_irc.Message) error {
	repo := plugin.repo
	return CommonRun(repo, plugin.ircClient, AddCommandPluginType, repo.AddCommand, command,
		channel, sender, message)
}

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
