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
	repo           repository.SingleBotRepositoryT
	messageHandler chathandler.ChatMessageHandlerT
}

var _ TwitchChatBotT = (*TwitchChatBot)(nil)

func NewTwitchChatBot(
	botInfo *models.Bot, ircClient client.TwitchClientT, botRepo repository.SingleBotRepositoryT,
	messageHandler chathandler.ChatMessageHandlerT) *TwitchChatBot {
	return &TwitchChatBot{
		botInfo:        botInfo,
		repo:           botRepo,
		ircClient:      ircClient,
		messageHandler: messageHandler,
	}
}

func DefaultTwitchChatBot(botInfo *models.Bot, baseRepo repository.BaseRepositoryT) TwitchChatBotT {
	ircClient := client.NewTwitchClient(botInfo.Username, botInfo.OauthToken)
	botRepo := repository.NewSingleBotRepository(botInfo, baseRepo)
	handler := chathandler.DefaultChatMessageHandler(botInfo, botRepo, ircClient)
	return NewTwitchChatBot(botInfo, ircClient, botRepo, handler)
}

func (bot *TwitchChatBot) Start() {
	client := bot.ircClient
	client.OnPrivateMessage(bot.messageHandler.OnPrivateMessage)

	// Join all channels associated with this bot
	channelModels := bot.repo.GetAllChannels()
	channelNames := make([]string, len(channelModels))
	for i, channelModel := range channelModels {
		channelNames[i] = channelModel.Username
	}
	client.Join(channelNames...)

	err := client.Connect()
	if err != nil {
		log.Fatal(err)
	}
}

func (bot *TwitchChatBot) Shutdown() {
	client := bot.ircClient

	// TODO: Is it ok to not store connected channels in memory?
	channels := bot.repo.GetAllChannels()
	for _, channel := range channels {
		client.Depart(channel.Username)
	}

	client.Disconnect()

	bot.messageHandler.Close()
}
