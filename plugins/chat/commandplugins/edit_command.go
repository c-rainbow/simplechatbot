package commandplugins

import (
	bot "github.com/c-rainbow/simplechatbot"
	chat_plugins "github.com/c-rainbow/simplechatbot/plugins/chat"
	twitch_irc "github.com/gempir/go-twitch-irc"
)

const (
	EditCommandPluginType = "EditCommandPluginType"
)

type EditCommandPlugin struct {
	ircClient *bot.TwitchClient
	repo      bot.SingleBotRepositoryT
}

var _ chat_plugins.ChatCommandPlugin = (*EditCommandPlugin)(nil)

func (plugin *EditCommandPlugin) Run(
	commandName string, channel string, sender *twitch_irc.User, message *twitch_irc.Message) error {
	repo := plugin.repo
	return CommonRun(repo, plugin.ircClient, EditCommandPluginType, repo.EditCommand, commandName,
		channel, sender, message)
}
