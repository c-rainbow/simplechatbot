package simplechatbot

import (
	"flag"
	"log"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/c-rainbow/simplechatbot/models"
	"github.com/guregu/dynamo"
)

var DatabaseEndpoint = *flag.String("dynamodb-endpoint", "http://localhost:8000", "DynamoDB endpoint address")
var DatabaseRegion = *flag.String("dynamodb-region", "us-west-2", "Default Region for DynamoDB")
var DisableSSL = *flag.Bool("dynamodb-disable-ssl", true, "If true, disable SSL to connect to DynamoDB")

type CommandMapKey struct {
	botID   int64
	channel string
}

// RepositoryInterface interface to deal with persistent data
type RepositoryInterface interface {
	GetAllBots() []*models.Bot
	GetAllChannels() []*models.Channel

	// For handlers to find command
	GetCommandByChannelAndName(channel string, commandName string)
}

type BaseRepository struct {
	commandMap map[CommandMapKey]map[string]*models.Command
	db         *dynamo.DB
}

func NewBaseRepository() *BaseRepository {
	// Initialize flag values
	db := dynamo.New(session.New(), &aws.Config{
		Endpoint:   aws.String(DatabaseEndpoint),
		Region:     aws.String(DatabaseRegion),
		DisableSSL: aws.Bool(DisableSSL),
	})
	return &BaseRepository{db: db}
}

// GetAllBots returns all Bot models in the database.
func (repo *BaseRepository) GetAllBots() []*models.Bot {
	bots := []*models.Bot{}
	err := repo.db.Table("Bots").Scan().All(&bots)
	if err != nil {
		log.Fatal("Error while fetching all bots", err.Error())
	}
	return bots
}

// GetAllChannels returns all Channel models in the database.
func (repo *BaseRepository) GetAllChannels() []*models.Channel {
	channels := []*models.Channel{}
	err := repo.db.Table("Channels").Scan().All(&channels)
	if err != nil {
		log.Fatal("Error while fetching all channels", err.Error())
	}
	return channels
}

// GetCommands returns map of command name to command model, for the bot and the channel
// Empty map is returned if commands do not exist for the combination.
func (repo *BaseRepository) GetCommands(botID int64, channel string) map[string]*models.Command {
	key := CommandMapKey{botID, channel}
	return repo.commandMap[key]
}

// GetCommand gets command by unique combination of (botID, channel, commandName)
// returns nil if command does not exist with the combination
func (repo *BaseRepository) GetCommand(botID int64, channel string, commandName string) *models.Command {
	key := CommandMapKey{botID, channel}
	return repo.commandMap[key][commandName]
}
