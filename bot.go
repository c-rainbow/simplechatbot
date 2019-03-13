// Bot struct and functions

package simplechatbot

import (
	"log"

	models "github.com/c-rainbow/simplechatbot/models"
)

// TwitchChatBot Twitch chat bot struct
type TwitchChatBot struct {
	botInfo        *models.Bot
	ircClient      *TwitchClient
	repo           BaseRepositoryT
	messageHandler ChatMessageHandlerT
}

func NewTwitchChatBot(
	botInfo *models.Bot, ircClient *TwitchClient, repo BaseRepositoryT,
	messageHandler ChatMessageHandlerT) *TwitchChatBot {
	return &TwitchChatBot{
		ircClient:      ircClient,
		botInfo:        botInfo,
		repo:           repo,
		messageHandler: messageHandler,
	}
}

func (bot *TwitchChatBot) Connect() {
	client := bot.ircClient
	client.OnNewMessage(bot.messageHandler.OnNewMessage)

	// Join all channels associated with this bot
	channels := bot.repo.GetAllChannelsForBot(bot.botInfo.TwitchID)
	for _, channel := range channels {
		client.Join(channel.Username)
	}

	err := client.Connect()
	if err != nil {
		log.Fatal(err)
	}
}

func (bot *TwitchChatBot) Disconnect() {
	bot.ircClient.Disconnect()
}
