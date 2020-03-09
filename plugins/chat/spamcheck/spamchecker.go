package spamcheck

import (
	"github.com/c-rainbow/simplechatbot/client"
	chatplugins "github.com/c-rainbow/simplechatbot/plugins/chat"
	twitch_irc "github.com/gempir/go-twitch-irc"
)

/*
Plugin to check for spam messages in chat
It listens to all messages, and if a chat message looks like a chat, it either
(1) Warns the chatter
(2) Timeout the chatter for N seconds (N can be different, depends on previous timeout history)
    i) Also warns the chatter
(3) Permanently ban the chatter
    i) Also say something cool to the banned chatter
*/

const (
	// SpamCheckerPluginType plugin type name to check if a chat message is a spam
	SpamCheckerPluginType = "SpamCheckerPluginType"
)

// SpamCheckerPlugin plugin to to check if a chat message is a spam
type SpamCheckerPlugin struct {
	ircClient client.TwitchClientT
}

var _ chatplugins.ChatListenerPluginT = (*SpamCheckerPlugin)(nil)

// NewSpamCheckerPlugin creates a new SpamCheckerPlugin
func NewSpamCheckerPlugin(ircClient client.TwitchClientT) chatplugins.ChatListenerPluginT {
	return &SpamCheckerPlugin{ircClient: ircClient}
}

// GetPluginType return plugin type
func (plugin *SpamCheckerPlugin) GetPluginType() string {
	return SpamCheckerPluginType
}

// ListenToChat listens to chat
func (plugin *SpamCheckerPlugin) ListenToChat(channel string, sender *twitch_irc.User,
	message *twitch_irc.PrivateMessage) bool {
	/*if len(message.Message) > 20 {
		plugin.ircClient.Say(channel, "message too long")
		return false
	}*/
	return true
}
