// Twitch IRC client module

package client

import (
	twitch_irc "github.com/gempir/go-twitch-irc"
)

// General interface for Twitch IRC client
// TOOD: OnNewMessage is currently tied to go-twitch-irc library.
type TwitchClientT interface {
	Connect() error
	Disconnect() error

	Join(channel string)
	Depart(channel string)

	Say(channel, text string)

	OnConnect(callback func())
	OnNewMessage(func(channel string, sender twitch_irc.User, message twitch_irc.Message))
}

// TwitchClient is Wrapper struct for existing Twitch IRC client.
type TwitchClient struct {
	// Inner twitch IRC chat client
	*twitch_irc.Client
}

var _ TwitchClientT = (*TwitchClient)(nil)

// NewTwitchClient creates new client.
func NewTwitchClient(username, oauthToken string) TwitchClientT {
	return &TwitchClient{
		twitch_irc.NewClient(username, oauthToken),
	}
}
