package main

import (
	"fmt"
	"time"

	"github.com/c-rainbow/simplechatbot/api/helix"
	"github.com/c-rainbow/simplechatbot/bot"
	"github.com/c-rainbow/simplechatbot/botmanager"
	"github.com/c-rainbow/simplechatbot/botserver"
	"github.com/c-rainbow/simplechatbot/chathandler"
	"github.com/c-rainbow/simplechatbot/client"
	"github.com/c-rainbow/simplechatbot/db/localrun"
	"github.com/c-rainbow/simplechatbot/flags"
	"github.com/c-rainbow/simplechatbot/install"
	"github.com/c-rainbow/simplechatbot/models"
	"github.com/c-rainbow/simplechatbot/parser"
	pluginmanager "github.com/c-rainbow/simplechatbot/plugins/chat/manager"

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
	localrun.ResetChannelsTable()
}

func main11() {
	client := helix.DefaultHelixClient()
	users, err := client.GetUsers(nil, []string{"c_rainbow"})
	_ = err
	for _, user := range users {
		fmt.Printf("ID: %s Name: %s\n", user.ID, user.DisplayName)
		fmt.Printf("Description: %s Email: %s\n", user.Description, user.Email)
		// fmt.Printf("ID: %s Name: %s\n", user.ID, user.DisplayName)
	}
}

func mainInstall() {

	installer := install.NewInstallerKor()
	fmt.Println("Install start")
	err := installer.Install()
	fmt.Println("Install end")
	if err != nil {
		fmt.Println("error: ", err.Error())
	}
}

func main() {
	/*flags.ParseAllFlags()
	localrun.DeleteAllTables()
	mainInstall()

	if 1 == 1 {
		return
	}*/

	flags.ParseAllFlags()

	baseRepo := repository.NewBaseRepository()
	manager := botmanager.NewBotManagerFromRepository(baseRepo)
	server := botserver.NewBotServer(nil, manager)
	server.Start()
}

// Run bot
func main30() {

	flags.ParseAllFlags()
	/*if true {
		fmt.Println("Binary exe test...")
		//time.Sleep(10 * time.Second)
		return
	}*/

	if flags.ResetChannels {
		fmt.Println("Resetting channels...")
		localrun.ResetChannelsTable()
		return
	}

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
	pluginManager := pluginmanager.NewChatCommandPluginManager(ircClient, botRepo)
	handler := chathandler.NewChatMessageHandler(botModels[0], botRepo, ircClient, pluginManager)
	fmt.Println("line 5")
	bot := bot.NewTwitchChatBot(botModels[0], ircClient, baseRepo, handler)
	fmt.Println("line 6")
	bot.Start()
	fmt.Println("line 7")
	// chatBot.Disconnect()
	// fmt.Println("Disconnected in main.go")

}

func mainAddCommand() {
	baseRepo := repository.NewBaseRepository()

	responseMap := make(map[string]models.ParsedResponse)
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
