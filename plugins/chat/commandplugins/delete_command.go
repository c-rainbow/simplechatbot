package commandplugins

import (
	"github.com/c-rainbow/simplechatbot/client"
	chat_plugins "github.com/c-rainbow/simplechatbot/plugins/chat"
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

var _ chat_plugins.ChatCommandPlugin = (*DeleteCommandPlugin)(nil)

func (plugin *DeleteCommandPlugin) Run(
	commandName string, channel string, sender *twitch_irc.User, message *twitch_irc.Message) error {
	repo := plugin.repo
	return CommonRun(repo, plugin.ircClient, DeleteCommandPluginType, repo.DeleteCommand, commandName,
		channel, sender, message)
}