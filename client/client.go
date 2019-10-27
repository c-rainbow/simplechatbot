// Twitch IRC client module

package client

import (
	twitch_irc "github.com/gempir/go-twitch-irc"
)

// General interface for Twitch IRC client
// TOOD: Function signature of OnPrivateMessage is currently tied to go-twitch-irc library.
type TwitchClientT interface {
	Connect() error
	Disconnect() error

	Join(channel string)
	Depart(channel string)

	Say(channel, text string)

	OnConnect(callback func())
	OnPrivateMessage(callback func(message twitch_irc.PrivateMessage))
}

// TwitchClient is Wrapper struct for existing Twitch IRC client.
type TwitchClient struct {
	// Inner twitch IRC chat client
	ircClient *twitch_irc.Client
}

var _ TwitchClientT = (*TwitchClient)(nil)

// NewTwitchClient creates new client.
func NewTwitchClient(username, oauthToken string) TwitchClientT {
	return &TwitchClient{
		ircClient: twitch_irc.NewClient(username, oauthToken),
	}
}

func (client *TwitchClient) Connect() error {
	return client.ircClient.Connect()
}

func (client *TwitchClient) Disconnect() error {
	return client.ircClient.Disconnect()
}

func (client *TwitchClient) Join(channel string) {
	client.ircClient.Join(channel)
}

func (client *TwitchClient) Depart(channel string) {
	client.ircClient.Depart(channel)
}

func (client *TwitchClient) Say(channel, text string) {
	client.ircClient.Say(channel, text)
}

func (client *TwitchClient) OnConnect(callback func()) {
	client.ircClient.OnConnect(callback)
}

func (client *TwitchClient) OnPrivateMessage(callback func(message twitch_irc.PrivateMessage)) {
	client.ircClient.OnPrivateMessage(callback)
}
