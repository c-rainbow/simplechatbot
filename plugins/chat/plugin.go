package chatplugins

import (
	"github.com/c-rainbow/simplechatbot/models"
	twitch_irc "github.com/gempir/go-twitch-irc"
)

/*
All the base interfaces for chat-triggered plugins

Note that a plugin can be a type of both ChatCommandPluginT and ChatListenerPluginT.
For example, LastChatChecker plugin usually silently runs in background and update chatter's last chat timestamps,
but when called in chat by a command, for example "!lastseen [username]", it answers to the chat with the stored data.
*/

// ChatCommandPluginT is the common interface for plugins that are triggered by command in chat messages
type ChatCommandPluginT interface {
	GetPluginType() string
	ReactToChat(command *models.Command, channel string, sender *twitch_irc.User, message *twitch_irc.PrivateMessage)
}

// ChatListenerPluginT is the common interface for plugins that listen to all chats
type ChatListenerPluginT interface {
	GetPluginType() string
	ListenToChat(channel string, sender *twitch_irc.User, message *twitch_irc.PrivateMessage) bool
}
