package commandplugins

import (
	"github.com/c-rainbow/simplechatbot/client"
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
	commandName string, channel string, sender *twitch_irc.User, message *twitch_irc.Message) error {
	repo := plugin.repo
	return CommonRun(repo, plugin.ircClient, AddCommandPluginType, repo.AddCommand, commandName,
		channel, sender, message)
}

type AddCommandPluginFactory struct {
	ircClient client.TwitchClientT
	repo      repository.SingleBotRepositoryT
}

var _ chatplugins.ChatCommandPluginFactoryT = (*AddCommandPluginFactory)(nil)

func NewAddCommandPluginFactory(
	ircClient client.TwitchClientT, repo repository.SingleBotRepositoryT) chatplugins.ChatCommandPluginFactoryT {
	return &AddCommandPluginFactory{ircClient: ircClient, repo: repo}
}

func (plugin *AddCommandPluginFactory) GetPluginType() string {
	return AddCommandPluginType
}

func (plugin *AddCommandPluginFactory) BuildNewPlugin() chatplugins.ChatCommandPluginT {
	return NewAddCommandPlugin(plugin.ircClient, plugin.repo)
}
