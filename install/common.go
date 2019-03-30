package install

import (
	"errors"
	"log"
	"strconv"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/c-rainbow/simplechatbot/api/helix"
	"github.com/c-rainbow/simplechatbot/client"
	"github.com/c-rainbow/simplechatbot/config"
	"github.com/c-rainbow/simplechatbot/models"
	"github.com/c-rainbow/simplechatbot/repository"
	"github.com/go-ini/ini"
	"github.com/guregu/dynamo"
)

const (
	DefaultBotSection     = "DefaultBot"
	DefaultChannelSection = "DefaultChannel"
)

var (
	ErrNoHelixClient = errors.New("Cannot create a new Helix API client")
)

type InstallerMessages struct {
	NoHelixClient                string
	TwitchUsersAPIError          string
	TwitchBotAccountNotFound     string
	TwitchChannelAccountNotFound string
	ChatServerSuccessfulLogin    string
	ChatServerFailedLogin        string
	DynamoDBConnectionError      string
}

type DynamoDBConfig struct {
	endpoint   string
	region     string
	disableSSL bool
}

type Installer struct {
	iniFile         *ini.File
	ircClientFunc   func(string, string) client.TwitchClientT
	helixClientFunc func(string) helix.HelixClientT
	messages        *InstallerMessages
}

func (installer *Installer) Install() error {
	// Read all config values from INI file
	iniFile := installer.iniFile
	bot := ReadBot(iniFile)
	channel := ReadChannel(iniFile)
	dbConfig := ReadDynamoDB(iniFile)
	clientID := ReadClientID(iniFile)

	var db *dynamo.DB
	var helixClient helix.HelixClientT
	var baseRepo repository.BaseRepositoryT

	// First of all, check if bot username-oauthtoken combination works
	var err error
	// err := installer.TryAccessingChatServer(bot)
	if err == nil {
		// Create Helix API client
		helixClient, err = installer.CreateHelixClient(clientID)
	}
	if err == nil {
		err = installer.PopulateBotModel(helixClient, bot)
	}
	if err == nil {
		err = installer.PopulateChannelModel(helixClient, channel)
	}
	if err == nil {
		db, err = installer.TryAccessingDynamoDB(dbConfig)
	}
	if true {
		return nil
	}
	if err == nil {
		// Create Bots table
		err = db.CreateTable(repository.BotTableName, models.Bot{}).Run()
	}
	if err == nil {
		// Create Channels table
		err = db.CreateTable(repository.ChannelTableName, models.Channel{}).Run()
	}
	if err == nil {
		// Add Bot to Bots table
		baseRepo = repository.NewBaseRepositoryCustomDB(db)
		err = baseRepo.CreateNewBot(bot)
	}
	if err == nil {
		// Add channel to Channels table
		err = baseRepo.CreateNewChannel(channel)
	}
	if err == nil {
		// Add bot to channel
		err = baseRepo.AddBotToChannel(bot, channel)
	}
	return err
}

func (installer *Installer) CreateHelixClient(clientID string) (helix.HelixClientT, error) {
	helixClient := installer.helixClientFunc(clientID)
	// Check if Helix client is not nil
	if helixClient == nil {
		log.Println(installer.messages.NoHelixClient)
		return nil, ErrNoHelixClient
	}
	// TODO: how to check if helix client is working?
	return helixClient, nil
}

// Populate other fields of Bot model from config, by using Twitch Users API.
func (installer *Installer) PopulateBotModel(
	helixClient helix.HelixClientT, bot *models.Bot) error {
	users, err := helixClient.GetUsers(nil, []string{bot.Username})
	if err != nil {
		log.Println(installer.messages.TwitchUsersAPIError, err.Error())
		return err
	}
	if len(users) == 0 {
		log.Println(installer.messages.TwitchBotAccountNotFound)
		return ErrTwitchAccoutNotFound
	}

	botAccount := users[0]
	botID, _ := strconv.ParseInt(botAccount.ID, 10, 64)
	bot.TwitchID = botID
	return nil
}

func (installer *Installer) PopulateChannelModel(
	helixClient helix.HelixClientT, channel *models.Channel) error {
	users, err := helixClient.GetUsers(nil, []string{channel.Username})
	if err != nil {
		log.Println(installer.messages.TwitchUsersAPIError, err.Error())
		return err
	}
	if len(users) == 0 {
		log.Println(installer.messages.TwitchChannelAccountNotFound)
		return ErrTwitchAccoutNotFound
	}

	channelAccount := users[0]
	channelID, _ := strconv.ParseInt(channelAccount.ID, 10, 64)
	channel.TwitchID = channelID
	channel.DisplayName = channelAccount.DisplayName
	return nil
}

func (installer *Installer) TryAccessingChatServer(bot *models.Bot) error {
	// TODO: Consider using goroutine to prevent blocking
	ircClient := installer.ircClientFunc(bot.Username, bot.OauthToken)

	ircClient.OnConnect(func() {
		log.Println(installer.messages.ChatServerSuccessfulLogin)
		time.Sleep(3 * time.Second)
		ircClient.Disconnect()
	})

	err := ircClient.Connect()
	if err != nil {
		log.Println(installer.messages.ChatServerFailedLogin, err.Error())
	}
	return err
}

func (installer *Installer) TryAccessingDynamoDB(dbConfig *DynamoDBConfig) (*dynamo.DB, error) {
	db := dynamo.New(session.New(), &aws.Config{
		Endpoint:   aws.String(dbConfig.endpoint),
		Region:     aws.String(dbConfig.region),
		DisableSSL: aws.Bool(dbConfig.disableSSL),
	})

	// List tables to check connection
	_, err := db.ListTables().All()
	if err != nil {
		log.Println(installer.messages.DynamoDBConnectionError, err.Error())
	}
	return db, err
}

// Read bot model from INI file. The data here is not yet validated.
func ReadBot(iniFile *ini.File) *models.Bot {
	section := iniFile.Section(DefaultBotSection)
	username := section.Key("BotUsername").String()
	oauthToken := section.Key("BotOauthToken").String()
	return &models.Bot{Username: username, OauthToken: oauthToken}
}

// Read channel model from INI file. The data here is not yet validated.
func ReadChannel(iniFile *ini.File) *models.Channel {
	section := iniFile.Section(DefaultChannelSection)
	username := section.Key("ChannelUsername").String()
	return &models.Channel{Username: username}
}

func ReadDynamoDB(iniFile *ini.File) *DynamoDBConfig {
	dbSection := iniFile.Section(config.DynamoDBSection)
	endpoint := dbSection.Key("DynamoDBAddress").String()
	region := dbSection.Key("DynamoDBRegion").String()
	disableSSL := dbSection.Key("DynamoDisableSSL").MustBool(true)
	return &DynamoDBConfig{endpoint: endpoint, region: region, disableSSL: disableSSL}
}

func ReadClientID(iniFile *ini.File) string {
	apiSection := iniFile.Section(config.TwitchAPISection)
	return apiSection.Key("ClientID").String()
}
