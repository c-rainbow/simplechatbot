package manager

import (
	"github.com/c-rainbow/simplechatbot/models"
	twitch_irc "github.com/gempir/go-twitch-irc"
)

const (
	jobChanDefaultCapacity = 100
)

// WorkItem is an item in queue for plugin to process
type WorkItem struct {
	command *models.Command
	channel string
	sender  *twitch_irc.User
	message *twitch_irc.PrivateMessage
}
