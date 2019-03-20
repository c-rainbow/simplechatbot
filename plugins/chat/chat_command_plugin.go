package chatplugins

import (
	twitch_irc "github.com/gempir/go-twitch-irc"
)

/*
Common infrastructure for plugins that are triggered by command in chat messages
2019-03-19: All plugins should have a command, even if they are dummy no-op ones.
This is for future-compatibility, when plugins might be controlled by chat commands.

Note that The same plugin can be accessed by multiple paths. For example, LastChatChecker
plugin usually silently runs in background and update chatter's last chat timestamps,
but when called in chat by a command, such as "!lastseen [username]", then it answers
back to the chat with the stored data.
*/
type ChatCommandPlugin interface {
	Run(commandName string, channel string, sender *twitch_irc.User, message *twitch_irc.Message) error
}
