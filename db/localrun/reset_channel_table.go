package localrun

import (
	"fmt"
	"log"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/c-rainbow/simplechatbot/flags"
	"github.com/guregu/dynamo"

	"github.com/c-rainbow/simplechatbot/models"
	"github.com/c-rainbow/simplechatbot/parser"
	"github.com/c-rainbow/simplechatbot/plugins/chat/commandplugins"
	"github.com/c-rainbow/simplechatbot/plugins/chat/games"
	"github.com/c-rainbow/simplechatbot/plugins/chat/selfban"
	"github.com/c-rainbow/simplechatbot/repository"
)

// Delete and re-create Channels table
// Add default channel with basic commands.
func ResetChannelsTable() {
	var err error

	// First, recreate Channels table
	err = RecreateChannelsTable()
	if err != nil {
		log.Println("Error while recreating table:", err.Error())
		return
	}
	// Build Channel struct
	commandMap := make(map[string]models.Command)
	botInfo := &models.Bot{TwitchID: DefaultBotTwitchID}
	chanInfo := &models.Channel{
		TwitchID: DefaultChannelTwitchID,
		Username: DefaultChannelUsername,
		BotIDs:   []int64{botInfo.TwitchID},
		Commands: make(map[string]models.Command),
	}

	repo := repository.NewBaseRepository()
	botID := botInfo.TwitchID
	channelID := chanInfo.TwitchID
	commandMap = chanInfo.Commands

	// These are default commands a new channel should have

	// Commands to add new command
	AddCommandToMap(commandMap, botID, channelID, "!addcom", commandplugins.AddCommandPluginType,
		"@$(user) Successfully added a new command")
	AddCommandToMap(commandMap, botID, channelID, "!추가", commandplugins.AddCommandPluginType,
		"@$(user) 명령어를 성공적으로 추가하였습니다")

	// Commands to edit existing command
	AddCommandToMap(commandMap, botID, channelID, "!editcom", commandplugins.EditCommandPluginType,
		"@$(user) Successfully edited the command")
	AddCommandToMap(commandMap, botID, channelID, "!수정", commandplugins.EditCommandPluginType,
		"@$(user) 명령어를 성공적으로 수정하였습니다")

	// Commands to delete existing command
	AddCommandToMap(commandMap, botID, channelID, "!delcom", commandplugins.DeleteCommandPluginType,
		"@$(user) Successfully deleted the command")
	AddCommandToMap(commandMap, botID, channelID, "!삭제", commandplugins.DeleteCommandPluginType,
		"@$(user) 명령어를 성공적으로 삭제하였습니다")

	// Commands to list available commands
	AddCommandToMap(commandMap, botID, channelID, "!commands", commandplugins.ListCommandsPluginType,
		"The list of commands are: $(arg0)")
	AddCommandToMap(commandMap, botID, channelID, "!명령어", commandplugins.ListCommandsPluginType,
		"명령어 모음: $(arg0)")

	AddCommandToMap(commandMap, botID, channelID, "!숫자", games.NumberGuesserPluginType, "")

	AddCommandToMap(commandMap, botID, channelID, "!셀프밴", selfban.BanOneselfPluginType, "@$(user) 님 밴")
	AddCommandToMap(commandMap, botID, channelID, "!banme", selfban.BanOneselfPluginType,
		"@$(user) You are banned for 5 seconds")

	AddCommandToMap(commandMap, botID, channelID, "!주사위", games.DicePluginType, ".")

	err = repo.CreateNewChannel(chanInfo)
	if err != nil {
		fmt.Println("Error while adding channel ", err.Error())
	}
}

func AddCommandToMap(commandMap map[string]models.Command, botID int64, channelID int64, name string, pluginType string, defaultResponse string) {
	command := buildCommand(botID, channelID, name, pluginType, defaultResponse)
	commandMap[name] = *command
}

func RecreateChannelsTable() error {
	db := dynamo.New(session.New(), &aws.Config{
		Endpoint:   aws.String(flags.DatabaseEndpoint),
		Region:     aws.String(flags.DatabaseRegion),
		DisableSSL: aws.Bool(flags.DisableSSL),
	})
	var err error
	// Delete Channels table
	err = db.Table(repository.ChannelTableName).DeleteTable().Run()
	if err != nil {
		fmt.Println("Error while deleting Channels table. ", err.Error())
		//return err
	}

	time.Sleep(1 * time.Second)

	// Create Channels table
	err = db.
		CreateTable(repository.ChannelTableName, &models.Channel{}).     // Create table
		Project("Username-index", dynamo.IncludeProjection, "Username"). // Create index
		Run()
	if err != nil {
		fmt.Println("Error while creating Channels table: " + err.Error())
	}

	time.Sleep(1 * time.Second)

	return err
}

func buildCommand(
	botID int64, channelID int64, name string, pluginType string, defaultResponse string) *models.Command {

	responseMap := make(map[string]models.ParsedResponse)
	responseMap[models.DefaultResponseKey] = *parser.ParseResponse(defaultResponse)

	return &models.Command{
		BotID:          botID,
		ChannelID:      channelID,
		Name:           name,
		PluginType:     pluginType,
		Responses:      responseMap,
		CooldownSecond: 5,
		Permission:     models.PermissionEveryone,
		Enabled:        true,
		Group:          "",
	}
}
