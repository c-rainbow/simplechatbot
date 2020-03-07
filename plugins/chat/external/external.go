package external

import (
	"github.com/c-rainbow/simplechatbot/client"
	"github.com/c-rainbow/simplechatbot/models"
	chatplugins "github.com/c-rainbow/simplechatbot/plugins/chat"
	twitch_irc "github.com/gempir/go-twitch-irc"
)

const (
	ExternalCallPluginType = "ExternalCallPluginType"
)

// Plugin type that sends data to external service and gets response.
//

// TODO: Build protobuf from message, as minimal as possible
// TODO: Parse protobuf from response from external call.

// TODO: Build JSON from message, as minimal as possible
// TODO: Parse JSON from response from external call.

// Plugin that responds to user-defined chat commands.
type ExternalCallPlugin struct {
	ircClient client.TwitchClientT
}

var _ chatplugins.ChatCommandPluginT = (*ExternalCallPlugin)(nil)

func NewExternalCallPlugin(ircClient client.TwitchClientT) chatplugins.ChatCommandPluginT {
	return &ExternalCallPlugin{ircClient: ircClient}
}

func (plugin *ExternalCallPlugin) GetPluginType() string {
	return ExternalCallPluginType
}

func (plugin *ExternalCallPlugin) ReactToChat(
	command *models.Command, channel string, sender *twitch_irc.User, message *twitch_irc.PrivateMessage) {
}
