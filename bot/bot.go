// Bot struct and functions

package bot

import (
	"log"

	chathandler "github.com/c-rainbow/simplechatbot/chathandler"
	client "github.com/c-rainbow/simplechatbot/client"
	models "github.com/c-rainbow/simplechatbot/models"
	repository "github.com/c-rainbow/simplechatbot/repository"
)

type TwitchChatBotT interface {
	Start()
	Shutdown()
}

// TwitchChatBot Twitch chat bot struct
type TwitchChatBot struct {
	botInfo        *models.Bot
	ircClient      client.TwitchClientT
	repo           repository.BaseRepositoryT
	messageHandler chathandler.ChatMessageHandlerT
}

var _ TwitchChatBotT = (*TwitchChatBot)(nil)

func NewTwitchChatBot(
	botInfo *models.Bot, ircClient client.TwitchClientT, repo repository.BaseRepositoryT,
	messageHandler chathandler.ChatMessageHandlerT) *TwitchChatBot {
	return &TwitchChatBot{
		ircClient:      ircClient,
		botInfo:        botInfo,
		repo:           repo,
		messageHandler: messageHandler,
	}
}

func (bot *TwitchChatBot) Start() {
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

func (bot *TwitchChatBot) Shutdown() {
	// TODO: close chans
	bot.ircClient.Disconnect()
}
