package plugins

import (
	twitch_irc "github.com/gempir/go-twitch-irc"
)

// Common infrastructure for plugins that are triggered by chat messages
type ChatPlugin interface {
	Run(commandName string, channel string, sender *twitch_irc.User, message *twitch_irc.Message) error
}
