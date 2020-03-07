package lastchatcheck

import (
	"github.com/c-rainbow/simplechatbot/client"
	"github.com/c-rainbow/simplechatbot/repository"
)

/*
Plugin to check last chat timestamp of chatter
It listens to all chats, and silently updates the timestamp

On the other hand, this plugin can be also called from chat by a command, such as !lastseen $(user),
*/

// CommandResponsePlugin plugin that responds to user-defined chat commands.
type CommandResponsePlugin struct {
	ircClient client.TwitchClientT
	repo      repository.SingleBotRepositoryT
}

//var _ chatplugins.ChatCommandPluginT = (*CommandResponsePlugin)(nil)
