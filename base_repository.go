package simplechatbot

import (
	"flag"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/c-rainbow/simplechatbot/models"
	"github.com/guregu/dynamo"
)

var DatabaseEndpoint = *flag.String("dynamodb-endpoint", "", "DynamoDB endpoint address")
var DatabaseRegion = *flag.String("dynamodb-region", "us-west-2", "Default Region for DynamoDB")
var DisableSSL = *flag.Bool("dynamodb-disable-ssl", true, "If true, disable SSL to connect to DynamoDB")

// RepositoryInterface interface to deal with persistent data
type RepositoryInterface interface {
	GetAllBots() []*models.Bot
	GetAllChannels() []*models.Channel

	// For handlers to find command
	GetCommandByChannelAndName(channel string, commandName string)
}

type BaseRepository struct {
	db *dynamo.DB
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

func (repo *BaseRepository) GetCommandByChannelAndName(channel string, command string) *models.Command {
	return nil
}

func (repo *BaseRepository) GetAllBots() []*models.Bot {
	iter := repo.db.Table("Bots").Scan().Iter()
	bots := []*models.Bot{}
	for {
		bot := models.Bot{}
		hasNext := iter.Next(&bot)
		if !hasNext {
			break
		}
		bots = append(bots, &bot)
	}
	return bots
}

func (repo *BaseRepository) GetAllChannels() []*models.Channel {
	iter := repo.db.Table("Channels").Scan().Iter()
	channels := []*models.Channel{}
	for {
		channel := models.Channel{}
		hasNext := iter.Next(&channel)
		if !hasNext {
			break
		}
		channels = append(channels, &channel)
	}
	return channels
}

/*
func (repo *Repository) scanAllItems(tableName string) dynamo.PagingIter {
	table := repo.db.Table(tableName)z
	iter := table.Scan().Iter()
	return iter
}
*/
// var _ DataRepositoryInterface = (*DataRepository)(nil)

// func (repo *DataRepository) GetAllBots() []*models.Bot {
// 	return []*models.Bot{}
// }

// func (repo *DataRepository)
