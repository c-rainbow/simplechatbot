package manager

import (
	"github.com/c-rainbow/simplechatbot/client"
	chatplugins "github.com/c-rainbow/simplechatbot/plugins/chat"
	"github.com/c-rainbow/simplechatbot/plugins/chat/spamcheck"
	"github.com/c-rainbow/simplechatbot/repository"
	twitch_irc "github.com/gempir/go-twitch-irc"
)

const (
	defaultCapacity = 100
)

// ChatListenerPluginManagerT interface for chat listener plugin manager
type ChatListenerPluginManagerT interface {
	ProcessChat(channel string, sender *twitch_irc.User, message *twitch_irc.PrivateMessage) bool
}

/*
ChatListenerPluginManager Chat listener plugin manager. All chat listener plugins are registered here.

A chat listener plugin listens to all chats in the channel. They work with the following rules:

(1) Chats in the same Twitch channel are processed sequentially, in first-come, first-served manner.
(2) For each chat, plugins run in order that they are registered in the manager.

The order of registered plugins is important because a common usage of listener plugin is to punish
(ban, timeout, warn, etc) inappropriate chats. For example, a long emote-only chat can be punished for
(1) too long, or (2) spamming emotes. If one plugin decides to timeout the user, then some plugins should
not process the chat anymore.
*/
type ChatListenerPluginManager struct {
	plugins []chatplugins.ChatListenerPluginT
}

var _ ChatListenerPluginManagerT = (*ChatListenerPluginManager)(nil)

// NewChatListenerPluginManager creates a new ChatListenerPluginManager
func NewChatListenerPluginManager(
	ircClient client.TwitchClientT, repo repository.SingleBotRepositoryT) ChatListenerPluginManagerT {
	manager := ChatListenerPluginManager{plugins: []chatplugins.ChatListenerPluginT{}}

	manager.RegisterPlugin(spamcheck.NewSpamCheckerPlugin(ircClient))

	return &manager
}

// RegisterPlugin registers a new plugin to the manager.
func (manager *ChatListenerPluginManager) RegisterPlugin(plugin chatplugins.ChatListenerPluginT) {
	manager.plugins = append(manager.plugins, plugin)
}

// ProcessChat parses command form chat message and calls the right plugin if necessary.
func (manager *ChatListenerPluginManager) ProcessChat(
	channel string, sender *twitch_irc.User, message *twitch_irc.PrivateMessage) bool {
	toContinue := true
	for _, plugin := range manager.plugins {
		if toContinue {
			toContinue = plugin.ListenToChat(channel, sender, message)
		} else {
			break
		}
	}
	return toContinue
}
