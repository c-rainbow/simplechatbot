// Bot struct and functions

package simplechatbot

import (
	"fmt"
	"log"
)

// TwitchChatBot Twitch chat bot struct
type TwitchChatBot struct {
	ircClient      *TwitchClient
	botInfo        *Bot
	userInfos      []*User
	messageHandler *ChatMessageHandler
}

func NewDefaultChatBot(botInfo *Bot) *TwitchChatBot {
	ircClient := NewTwitchClient(botInfo.Name, botInfo.OauthToken)
	fmt.Println("bot info: " + botInfo.Name)

	var botDataManager BotDataManager = &PlainTextBotDataManager{}
	userInfos := botDataManager.GetUserInfos()
	messageHandler := NewChatMessageHandler(ircClient, &botDataManager)

	return &TwitchChatBot{
		ircClient:      ircClient,
		botInfo:        botInfo,
		userInfos:      userInfos,
		messageHandler: messageHandler,
	}
}

func (bot *TwitchChatBot) Disconnect() {
	fmt.Println("disconnecting.. 1")
	bot.ircClient.Disconnect()
	fmt.Println("disconnecting.. 2")
}

func (bot *TwitchChatBot) Connect() {
	fmt.Println("connected? 1")
	client := bot.ircClient
	fmt.Println("connected? 2")
	// TODO: add callback functions for operations
	// client.OnNewWhisper(func(user twitch.User, message twitch.Message) {})
	client.OnNewMessage(bot.messageHandler.OnNewMessage)
	fmt.Println("connected? 3")

	// client.OnNewRoomstateMessage(func(channel string, user twitch.User, message twitch.Message) {})
	// client.OnNewClearchatMessage(func(channel string, user twitch.User, message twitch.Message) {})
	// client.OnNewUsernoticeMessage(func(channel string, user twitch.User, message twitch.Message) {})
	// client.OnNewNoticeMessage(func(channel string, user twitch.User, message twitch.Message) {})
	// client.OnNewUserstateMessage(func(channel string, user twitch.User, message twitch.Message) {})
	// client.OnUserJoin(func(channel, user string) {})
	// client.OnUserPart(func(channel, user string) {})

	// Join all channels assigned to this bot
	for _, userInfo := range bot.userInfos {
		client.Join(userInfo.Username)
		fmt.Print("Joined " + userInfo.Username)
	}
	fmt.Println("connected? 4")
	err := client.Connect()
	fmt.Println("connected? 5")
	if err != nil {
		log.Fatal(err)
		panic(err)
	}
}
