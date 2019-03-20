package commandplugins

import (
	"github.com/c-rainbow/simplechatbot/client"
	chat_plugins "github.com/c-rainbow/simplechatbot/plugins/chat"
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

var _ chat_plugins.ChatCommandPlugin = (*AddCommandPlugin)(nil)

func NewAddCommandPlugin(
	ircClient client.TwitchClientT, repo repository.SingleBotRepositoryT) chat_plugins.ChatCommandPlugin {
	return &AddCommandPlugin{ircClient: ircClient, repo: repo}
}

func (plugin *AddCommandPlugin) Run(
	commandName string, channel string, sender *twitch_irc.User, message *twitch_irc.Message) error {
	repo := plugin.repo
	return CommonRun(repo, plugin.ircClient, AddCommandPluginType, repo.AddCommand, commandName,
		channel, sender, message)
}
