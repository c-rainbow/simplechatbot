package commandplugins

import (
	"github.com/c-rainbow/simplechatbot/client"
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
	commandName string, channel string, sender *twitch_irc.User, message *twitch_irc.Message) error {
	repo := plugin.repo
	return CommonRun(repo, plugin.ircClient, DeleteCommandPluginType, repo.DeleteCommand, commandName,
		channel, sender, message)
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
