package main

import (
	// "fmt"

	"fmt"
	"time"

	"github.com/c-rainbow/simplechatbot/bot"
	"github.com/c-rainbow/simplechatbot/chathandler"
	"github.com/c-rainbow/simplechatbot/client"
	"github.com/c-rainbow/simplechatbot/db/localrun"
	"github.com/c-rainbow/simplechatbot/models"
	"github.com/c-rainbow/simplechatbot/parser"
	repository "github.com/c-rainbow/simplechatbot/repository"
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

func main2() {
	localrun.AddNewCommand()
}

// Run bot
func main() {

	baseRepo := repository.NewBaseRepository()
	fmt.Println("line 1")
	// chanModels := baseRepo.GetAllChannels()
	// singleChanRepo := repository.NewSingleChannelRepository(chanModels[0], baseRepo)
	botModels := baseRepo.GetAllBots()
	fmt.Println("line 2")
	botRepo := repository.NewSingleBotRepository(botModels[0], baseRepo)
	fmt.Println("line 3")
	ircClient := client.NewTwitchClient(botModels[0].Username, botModels[0].OauthToken)
	fmt.Println("line 4")
	handler := chathandler.NewChatMessageHandler(botModels[0], botRepo, ircClient)
	fmt.Println("line 5")
	bot := bot.NewTwitchChatBot(botModels[0], ircClient, baseRepo, handler)
	fmt.Println("line 6")
	bot.Connect()
	fmt.Println("line 7")
	// chatBot.Disconnect()
	// fmt.Println("Disconnected in main.go")

}

func mainAddCommand() {
	baseRepo := repository.NewBaseRepository()

	responseMap := make(map[string]parser.ParsedResponse)
	response := parser.ParseResponse("Hello $(user)")
	responseMap[""] = *response
	channelName := localrun.DefaultChannelUsername
	helloCommand := models.Command{
		Name:           "hello",
		BotID:          localrun.DefaultBotTwitchID,
		ChannelID:      localrun.DefaultChannelTwitchID,
		CooldownSecond: 5,
		Responses:      responseMap,
		Enabled:        true,
	}
	err := baseRepo.AddCommand(channelName, &helloCommand)
	if err != nil {
		fmt.Println("Errorrrrr: ", err.Error())
	}
}

func mainDelete() {
	baseRepo := repository.NewBaseRepository()
	command := models.Command{
		Name:           "hello",
		BotID:          localrun.DefaultBotTwitchID,
		ChannelID:      localrun.DefaultChannelTwitchID,
		CooldownSecond: 3,
	}
	err := baseRepo.DeleteCommand(localrun.DefaultChannelUsername, &command)
	if err != nil {
		fmt.Println("DeleteCommand failed: ", err.Error())
	}
}
