// Bot struct and functions

package bot

import (
	"log"

	chathandler "github.com/c-rainbow/simplechatbot/chathandler"
	client "github.com/c-rainbow/simplechatbot/client"
	models "github.com/c-rainbow/simplechatbot/models"
	pluginmanager "github.com/c-rainbow/simplechatbot/plugins/chat/manager"
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

func NewTwitchChatBotFromRepsitory(botInfo *models.Bot, baseRepo repository.BaseRepositoryT) TwitchChatBotT {
	ircClient := client.NewTwitchClient(botInfo.Username, botInfo.OauthToken)
	botRepo := repository.NewSingleBotRepository(botInfo, baseRepo)
	pluginManager := pluginmanager.NewChatCommandPluginManager(ircClient, botRepo)
	handler := chathandler.NewChatMessageHandler(botInfo, botRepo, ircClient, pluginManager)
	return NewTwitchChatBot(botInfo, ircClient, baseRepo, handler)
}

func (bot *TwitchChatBot) Start() {
	client := bot.ircClient
	client.OnPrivateMessage(bot.messageHandler.OnPrivateMessage)

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
	client := bot.ircClient

	// TODO: Is it ok to not store connected channels in memory?
	channels := bot.repo.GetAllChannelsForBot(bot.botInfo.TwitchID)
	for _, channel := range channels {
		client.Depart(channel.Username)
	}

	client.Disconnect()
}
