package install

import (
	"errors"
	"log"
	"strconv"
	"time"

	"github.com/c-rainbow/simplechatbot/repository"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/guregu/dynamo"

	"github.com/c-rainbow/simplechatbot/api/helix"
	"github.com/c-rainbow/simplechatbot/client"
	"github.com/c-rainbow/simplechatbot/config"
	"github.com/c-rainbow/simplechatbot/models"
	"github.com/go-ini/ini"
)

var (
	ErrTwitchAccoutNotFound = errors.New("Twitch Account Not Found")
)

/*

각각의 단계에서 안되면 실패

1. 설정파일 읽어오고 확인하기. 필요한 값들 모두 있는지 최소한 형식 (int, string, etc) 은 맞는지.

2. 설정파일 봇 존재하는 계정인지 확인
3. 설정파일 봇으로 채팅서버에 로그인 해보기. 테스트 채널에서 채팅 쳐보기 (어떻게 잘 되었는지 확인?)

4. 설정파일 채널 존재하는 계정인지 확인

5. DynamoDB 접속해 보기
6. DB에 Bots 테이블 만들기
7. DB에 Channels 테이블 만들기
8. 기본 봇 데이타 넣기
9. 기본 채널 데이타 넣기
10. 채널에 봇 넣기

*/

type InstallerKor struct {
	iniFile     *ini.File
	ircClient   client.TwitchClientT
	helixClient helix.HelixClientT
}

func NewInstallerKor() *InstallerKor {
	iniFile := ReadConfig()

	botSection := iniFile.Section(DefaultBotSection)
	botUsername := botSection.Key("BotUsername").String()
	botOauthToken := botSection.Key("BotOauthToken").String()
	ircClient := client.NewTwitchClient(botUsername, botOauthToken)

	helixClient := helix.DefaultHelixClient()
	if iniFile == nil || helixClient == nil {
		log.Println("Cannot get new installer")
		return nil
	}
	return &InstallerKor{iniFile: iniFile, ircClient: ircClient, helixClient: helixClient}
}

func (installer *InstallerKor) Install() error {
	var bot *models.Bot
	var channel *models.Channel
	var db *dynamo.DB
	var baseRepo repository.BaseRepositoryT
	bot, err := installer.CheckIfBotAccountExists()
	if err == nil {
		err = installer.TryAccessingChatServer(bot)
	}
	if err == nil {
		channel, err = installer.CheckIfChannelAccountExists()
	}

	if err == nil {
		db, err = installer.TryAccessingDynamoDB()
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

func (installer *InstallerKor) CheckIfBotAccountExists() (*models.Bot, error) {
	botSection := installer.iniFile.Section(DefaultBotSection)
	botUsername := botSection.Key("BotUsername").String()
	botOauthToken := botSection.Key("BotOauthToken").String()

	users, err := installer.helixClient.GetUsers(nil, []string{botUsername})
	if err != nil {
		return nil, err
	}
	if len(users) == 0 {
		return nil, ErrTwitchAccoutNotFound
	}

	botAccount := users[0]
	botID, err := strconv.ParseInt(botAccount.ID, 10, 64)
	if err != nil {
		return nil, err // No way this will be error, just checking.
	}
	return &models.Bot{TwitchID: botID, Username: botAccount.Login, OauthToken: botOauthToken}, nil
}

func (installer *InstallerKor) CheckIfChannelAccountExists() (*models.Channel, error) {
	channelSection := installer.iniFile.Section(DefaultChannelSection)
	channelUsername := channelSection.Key("ChannelUsername").String()

	users, err := installer.helixClient.GetUsers(nil, []string{channelUsername})
	if err != nil {
		return nil, err
	}
	if len(users) == 0 {
		return nil, ErrTwitchAccoutNotFound
	}

	channelAccount := users[0]
	channelID, err := strconv.ParseInt(channelAccount.ID, 10, 64)
	if err != nil {
		return nil, err // No way this will be error, just checking.
	}
	return &models.Channel{TwitchID: channelID, Username: channelAccount.Login}, nil
}

func ReadConfig() *ini.File {
	iniFile, err := config.NewConfig()
	if err != nil {
		log.Println("Error while reading config file", err.Error())
		return nil
	}
	return iniFile
}

func (installer *InstallerKor) TryAccessingChatServer(bot *models.Bot) error {
	// TODO: Consider using goroutine to prevent blocking
	installer.ircClient.OnConnect(func() {
		log.Println("Successfully connected to Twitch chat server. Disconnecting in 3 seconds")
		time.Sleep(3 * time.Second)
		installer.ircClient.Disconnect()
	})

	err := installer.ircClient.Connect()
	return err
}

func (installer *InstallerKor) TryAccessingDynamoDB() (*dynamo.DB, error) {
	dbSection := installer.iniFile.Section(DynamoDBSection)
	dbAddress := dbSection.Key("DynamoDBAddress").String()
	dbRegion := dbSection.Key("DynamoDBRegion").String()
	dbDisableSSL := dbSection.Key("DynamoDisableSSL").MustBool(true)

	db := dynamo.New(session.New(), &aws.Config{
		Endpoint:   aws.String(dbAddress),
		Region:     aws.String(dbRegion),
		DisableSSL: aws.Bool(dbDisableSSL),
	})

	// List tables to check connection
	_, err := db.ListTables().All()
	return db, err
}
