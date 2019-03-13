package main

import (
	// "fmt"

	"fmt"
	"time"

	"github.com/c-rainbow/simplechatbot"
	"github.com/c-rainbow/simplechatbot/db/localrun"
)

// Clean and re-populate DynamoDB tables with default bot and channel data
// Fill in constants.go before running this function
// TODO: DynamoDB often does not recognize recently created/deleted table
func main1() {
	localrun.DeleteAllTables()
	time.Sleep(2 * time.Second)
	fmt.Println("check 1")

	localrun.CreateAllTables()
	time.Sleep(2 * time.Second)
	fmt.Println("check 2")

	localrun.AddNewBot()
	time.Sleep(2 * time.Second)
	fmt.Println("check 3")

	localrun.AddNewChannel()
	time.Sleep(2 * time.Second)
	fmt.Println("check 4")

	//localrun.GetChannelsForBot()
	// time.Sleep(2 * time.Second)
	localrun.AddBotToChannel()
	time.Sleep(2 * time.Second)
	fmt.Println("check 5")
}

// Run bot
func main() {
	baseRepo := simplechatbot.NewBaseRepository()
	fmt.Println("line 1")
	// chanModels := baseRepo.GetAllChannels()
	// singleChanRepo := simplechatbot.NewSingleChannelRepository(chanModels[0], baseRepo)
	botModels := baseRepo.GetAllBots()
	fmt.Println("line 2")
	botRepo := simplechatbot.NewSingleBotRepository(botModels[0], baseRepo)
	fmt.Println("line 3")
	ircClient := simplechatbot.NewTwitchClient(botModels[0].Username, botModels[0].OauthToken)
	fmt.Println("line 4")
	handler := simplechatbot.NewChatMessageHandler(botModels[0], botRepo, ircClient)
	fmt.Println("line 5")
	bot := simplechatbot.NewTwitchChatBot(botModels[0], ircClient, baseRepo, handler)
	fmt.Println("line 6")
	bot.Connect()
	fmt.Println("line 7")
	// chatBot.Disconnect()
	// fmt.Println("Disconnected in main.go")
}
