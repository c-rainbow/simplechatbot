package simplechatbot

import (
	"log"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	models "github.com/c-rainbow/simplechatbot/models"
	dynamo "github.com/guregu/dynamo"
)

type SingleBotRepositoryT interface {
	GetAllChannels() []*models.Channel
	GetCommandByChannelAndName(channel string, commandName string) *models.Command
}

// Repository for a single bot
type SingleBotRepository struct {
	botInfo  *models.Bot
	baseRepo BaseRepositoryT
	db       *dynamo.DB
}

var _ SingleBotRepositoryT = (*SingleBotRepository)(nil)

func NewSingleBotRepository(botInfo *models.Bot, baseRepo BaseRepositoryT) *SingleBotRepository {
	// Initialize flag values
	db := dynamo.New(session.New(), &aws.Config{
		Endpoint:   aws.String(DatabaseEndpoint),
		Region:     aws.String(DatabaseRegion),
		DisableSSL: aws.Bool(DisableSSL),
	})

	return &SingleBotRepository{botInfo: botInfo, db: db}
}

// GetAllChannels returns all channels for this bot.
func (repo *SingleBotRepository) GetAllChannels() []*models.Channel {
	botID := repo.botInfo.TwitchID
	channels := []*models.Channel{}
	err := repo.db.Table("Channels").Scan().Filter("contains(BotIDs, ?)", botID).All(&channels)
	if err != nil {
		log.Fatal("Error while finding channels for bot", err.Error())
	}

	return channels
}

func (repo *SingleBotRepository) GetCommandByChannelAndName(
	channel string, commandName string) *models.Command {
	return repo.baseRepo.GetCommand(repo.botInfo.TwitchID, channel, commandName)
}
