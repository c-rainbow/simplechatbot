package repository

import (
	"errors"
	"fmt"
	"log"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	flags "github.com/c-rainbow/simplechatbot/flags"
	"github.com/c-rainbow/simplechatbot/models"
	"github.com/guregu/dynamo"
)

const (
	BotTableName     = "Bots"
	ChannelTableName = "Channels"
)

var (
	CommandNameNotFoundError  = errors.New("Command name is not found")
	CommandAlreadyExistsError = errors.New("Command name already exists")
)

/*type CommandMapKey struct {
	botID   int64
	channel string
}*/

// BaseRepositoryT interface to deal with persistent data
type BaseRepositoryT interface {
	GetAllBots() []*models.Bot
	GetAllChannels() []*models.Channel
	GetAllChannelsForBot(botID int64) []*models.Channel
	/// GetCommands(botID int64, channel string) map[string]*models.Command
	GetCommand(botID int64, channel string, commandName string) *models.Command
	CreateNewBot(botInfo *models.Bot) error
	AddBotToChannel(botInfo *models.Bot, channelToAdd *models.Channel) error
	AddCommand(channel string, commandToAdd *models.Command) error
	EditCommand(channel string, commandToAdd *models.Command) error
	DeleteCommand(channel string, commandToAdd *models.Command) error
	ListCommands(channel string) ([]*models.Command, error)
}

type BaseRepository struct {
	channelMap map[string]*models.Channel // Channel name -> command name -> command model
	db         *dynamo.DB
}

var _ BaseRepositoryT = (*BaseRepository)(nil)

func NewBaseRepository() *BaseRepository {
	// Initialize flag values
	db := dynamo.New(session.New(), &aws.Config{
		Endpoint:   aws.String(flags.DatabaseEndpoint),
		Region:     aws.String(flags.DatabaseRegion),
		DisableSSL: aws.Bool(flags.DisableSSL),
	})

	channelMap := make(map[string]*models.Channel)

	repo := &BaseRepository{db: db, channelMap: channelMap}
	// This should be called after initialization because it uses another function of BaseRepository
	repo.populateChannelMap()
	return repo
}

func (repo *BaseRepository) populateChannelMap() {
	channels := repo.GetAllChannels()
	for _, channel := range channels {
		if channel.Commands == nil {
			channel.Commands = make(map[string]models.Command)
		}
		repo.channelMap[channel.Username] = channel
	}
}

// GetAllBots returns all Bot models in the database.
func (repo *BaseRepository) GetAllBots() []*models.Bot {
	bots := []*models.Bot{}
	err := repo.db.Table(BotTableName).Scan().All(&bots)
	if err != nil {
		log.Fatal("Error while fetching all bots", err.Error())
	}
	return bots
}

// GetAllChannels returns all Channel models in the database.
func (repo *BaseRepository) GetAllChannels() []*models.Channel {
	channels := []*models.Channel{}
	err := repo.db.Table(ChannelTableName).Scan().All(&channels)
	if err != nil {
		log.Fatal("Error while fetching all channels", err.Error())
	}
	return channels
}

// GetAllChannels returns all channels for this bot.
func (repo *BaseRepository) GetAllChannelsForBot(botID int64) []*models.Channel {
	channels := []*models.Channel{}
	err := repo.db.Table(ChannelTableName).Scan().Filter(
		"contains(BotIDs, ?)", botID).All(&channels)
	if err != nil {
		log.Fatal("Error while finding channels for bot", err.Error())
	}

	return channels
}

// GetCommands returns map of command name to command model, for the bot and the channel
// Empty map is returned if commands do not exist for the combination.
/*func (repo *BaseRepository) GetCommands(botID int64, channel string) map[string]models.Command {
	return repo.channelMap[channel].Commands
}*/

// GetCommand gets command by unique combination of (botID, channel, commandName)
// returns nil if command does not exist with the combination
func (repo *BaseRepository) GetCommand(botID int64, channel string, commandName string) *models.Command {
	command, exists := repo.channelMap[channel].Commands[commandName]
	if exists && command.BotID == botID {
		return &command
	} else {
		return nil
	}
}

// CreateNewBot adds a new bot to DynamoDB.
// The function assumes that bot with the same key (TwitchID) does not already exist
func (repo *BaseRepository) CreateNewBot(botInfo *models.Bot) error {
	botTable := repo.db.Table(BotTableName)
	err := botTable.Put(botInfo).Run()
	return err
}

// CreateNewChannel adds a new channel to DynamoDB.
// The function assumes that channel with the same key (TwitchID) does not already exist
func (repo *BaseRepository) CreateNewChannel(chanInfo *models.Channel) error {
	chanTable := repo.db.Table(ChannelTableName)
	err := chanTable.Put(chanInfo).Run()
	return err
}

// AddBotToChannel adds the bot to the channel.
// It assumes that the bot is not already running in the channel.
// TODO: needs testing
func (repo *BaseRepository) AddBotToChannel(botInfo *models.Bot, channelToAdd *models.Channel) error {
	// Procedures
	// 1. Check the bot ID is correct

	// 2. Check the channel name exists && channel info is up-to-date.
	chanTable := repo.db.Table(ChannelTableName)
	channel := &models.Channel{}
	err := chanTable.Get("ID", channelToAdd.TwitchID).Filter(
		"$ = ?", "Username", channelToAdd.Username).One(&channel)
	if err != nil {
		log.Fatal("Error while retrieving channelll: ", err.Error())
		return err
	}

	// 3. Check if the bot already runs in the channel
	botExists := false
	for _, botID := range channel.BotIDs {
		if botID == botInfo.TwitchID {
			botExists = true
			fmt.Println("Bot is already running in the channel")
			return nil
		}
	}
	// 4. Add the botID to the channel
	if !botExists {
		chanTable := repo.db.Table(ChannelTableName)
		channel.BotIDs = append(channel.BotIDs, botInfo.TwitchID)
		err := chanTable.Put(channel).Run()
		return err
	}
	return nil

	// 5. (Outside of this file) Bot joins the channel
}

func (repo *BaseRepository) UpdateChannel(chanInfo *models.Channel) error {
	chanTable := repo.db.Table(ChannelTableName)
	err := chanTable.Put(chanInfo).Run()
	return err
}

func (repo *BaseRepository) AddCommand(channel string, commandToAdd *models.Command) error {
	chanInfo, exists := repo.getChannelAndCommandExists(channel, commandToAdd)
	if exists {
		return CommandAlreadyExistsError
	}
	chanInfo.Commands[commandToAdd.Name] = *commandToAdd

	// TODO: Handle when DB update failed, then we have inconsistent state in channelMap
	return repo.UpdateChannel(chanInfo)
}

func (repo *BaseRepository) EditCommand(channel string, commandToEdit *models.Command) error {
	chanInfo, exists := repo.getChannelAndCommandExists(channel, commandToEdit)
	if !exists {
		return CommandNameNotFoundError
	}
	chanInfo.Commands[commandToEdit.Name] = *commandToEdit

	// TODO: Handle when DB update failed, then we have inconsistent state in channelMap
	return repo.UpdateChannel(chanInfo)
}

func (repo *BaseRepository) DeleteCommand(channel string, commandToDelete *models.Command) error {
	chanInfo, exists := repo.getChannelAndCommandExists(channel, commandToDelete)
	if !exists {
		return CommandNameNotFoundError
	}
	delete(chanInfo.Commands, commandToDelete.Name)
	if _, has := chanInfo.Commands[""]; has {
		delete(chanInfo.Commands, "")
	}

	// TODO: Handle when DB update failed, then we have inconsistent state in channelMap
	return repo.UpdateChannel(chanInfo)
}

// Returns channel info
func (repo *BaseRepository) getChannelAndCommandExists(
	channel string, command *models.Command) (*models.Channel, bool) {
	chanInfo := repo.channelMap[channel]
	_, exists := chanInfo.Commands[command.Name]
	return chanInfo, exists
}

func (repo *BaseRepository) ListCommands(channel string) ([]*models.Command, error) {
	// TODO: Which bot's commands should be returned?
	commandMap := repo.channelMap[channel].Commands
	commands := []models.Command{}
	for _, v := range commandMap {
		fmt.Println("Baserepo command name: ", v.Name)
		commands = append(commands, v)
		fmt.Println("Command so far: ", commands)
	}
	commandPointers := make([]*models.Command, len(commands))
	for i, _ := range commands {
		commandPointers[i] = &commands[i]
	}
	return commandPointers, nil
}
