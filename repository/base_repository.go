package repository

import (
	"errors"
	"log"
	"sync"

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
	ErrChannelNotFound        = errors.New("Channel is not found")
	ErrChannelAlreadyExists   = errors.New("Channel already exists")
	ErrBotAlreadyInChannel    = errors.New("Bot is already running in the channel")
)

// BaseRepositoryT interface to deal with persistent data
type BaseRepositoryT interface {
	GetAllBots() []*models.Bot
	GetAllChannels() []*models.Channel
	GetAllChannelsForBot(botID int64) []*models.Channel
	GetCommand(botID int64, channel string, commandName string) *models.Command
	CreateNewBot(botInfo *models.Bot) error
	CreateNewChannel(chanInfo *models.Channel) error
	AddBotToChannel(botInfo *models.Bot, channelToAdd *models.Channel) error
	AddCommand(channel string, commandToAdd *models.Command) error
	EditCommand(channel string, commandToEdit *models.Command) error
	DeleteCommand(channel string, commandToDelete *models.Command) error
	ListCommands(channel string) ([]*models.Command, error)
}

// Note: All public methods of this struct that reads/updates channelMap
// should start with mutex.RLock() or mutex.Lock().
// Private functions assume that locks are already obtained.
// TODO: Create mutex per channel, not per repository.
type BaseRepository struct {
	mutex      sync.RWMutex
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

	repo := &BaseRepository{db: db}
	// This should be called after initialization because it uses another function of BaseRepository
	repo.PopulateChannelMap()
	return repo
}

func (repo *BaseRepository) PopulateChannelMap() {
	repo.mutex.Lock()
	defer repo.mutex.Unlock()

	channelMap := make(map[string]*models.Channel)
	channels := repo.GetAllChannels()
	for _, channel := range channels {
		if channel.Commands == nil {
			channel.Commands = make(map[string]models.Command)
		}
		channelMap[channel.Username] = channel
	}
	repo.channelMap = channelMap
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

// Get Channel object from Twitch ID, or error if not found
func (repo *BaseRepository) getChannelFromTwitchID(channelID int64) (*models.Channel, error) {
	channel := &models.Channel{}
	err := repo.db.Table(ChannelTableName).Get("ID", channelID).One(channel)
	if err != nil {
		log.Println("Cannot get channel from DB. Error:", err.Error())
	}
	return channel, err
}

// Get Bot object from Twitch ID, or error if not found
func (repo *BaseRepository) getBotFromTwitchID(botID int64) (*models.Bot, error) {
	bot := &models.Bot{}
	err := repo.db.Table(BotTableName).Get("ID", botID).One(bot)
	if err != nil {
		log.Println("Cannot get Bot from DB. Error:", err.Error())
	}
	return bot, err
}

// GetCommand gets command by unique combination of (botID, channel, commandName)
// returns nil if command does not exist with the combination
func (repo *BaseRepository) GetCommand(botID int64, channel string, commandName string) *models.Command {
	repo.mutex.RLock()
	defer repo.mutex.RUnlock()

	// Check if channel exists.
	chanInfo := repo.channelMap[channel]
	if chanInfo == nil {
		return nil
	}
	// Check if the channel has command with the name and botID
	command, exists := chanInfo.Commands[commandName]
	if exists && command.BotID == botID {
		return &command
	}
	return nil
}

// CreateNewBot adds a new bot to DynamoDB.
// The function assumes that bot with the same key (TwitchID) does not already exist
func (repo *BaseRepository) CreateNewBot(botInfo *models.Bot) error {
	repo.mutex.Lock()
	defer repo.mutex.Unlock()
	_, err := repo.getBotFromTwitchID(botInfo.TwitchID)

	if err != nil {
		return err
	}
	botTable := repo.db.Table(BotTableName)
	return botTable.Put(botInfo).Run()
}

// CreateNewChannel adds a new channel to DynamoDB.
// The function assumes that channel with the same key (TwitchID) does not already exist
func (repo *BaseRepository) CreateNewChannel(channelToAdd *models.Channel) error {
	repo.mutex.Lock()
	defer repo.mutex.Unlock()

	// First, check that there is no user with the same Twitch ID.
	// Don't check for username because the streamer may call this right after changing username.
	_, err := repo.getChannelFromTwitchID(channelToAdd.TwitchID)
	if err != nil {
		return err
	}

	// updateChannel() adds the channel to DB and channelMap if not exists.
	return repo.updateChannel(channelToAdd)
}

// AddBotToChannel adds the bot to the channel.
// It assumes that the bot is not already running in the channel.
// TODO: needs testing
// TODO: This doesn't update channelMap
func (repo *BaseRepository) AddBotToChannel(botInfo *models.Bot, channelToAdd *models.Channel) error {
	repo.mutex.Lock()
	defer repo.mutex.Unlock()

	// Procedures
	// 1. Check the bot ID is correct

	// 2. Check that the channel exists
	channelName := channelToAdd.Username
	channel, exists := repo.channelMap[channelName]
	if !exists {
		return ErrChannelNotFound
	}
	chanInfo := *channel // Make a copy so that the original isn't modified until successful DB operation.

	// 3. Check if the bot already runs in the channel
	for _, botID := range chanInfo.BotIDs {
		if botID == botInfo.TwitchID {
			return ErrBotAlreadyInChannel
		}
	}

	// 4. Add the botID to the channel
	chanInfo.BotIDs = append(chanInfo.BotIDs, botInfo.TwitchID)
	return repo.updateChannel(&chanInfo)

	// 5. (Outside of this file) Bot joins the channel
}

// This function assumes that lock IS obtained by its caller
func (repo *BaseRepository) updateChannel(chanInfo *models.Channel) error {
	chanTable := repo.db.Table(ChannelTableName)
	err := chanTable.Put(chanInfo).Run()
	// Update channelMap only when DB operation is successful
	if err == nil {
		repo.channelMap[chanInfo.Username] = chanInfo
	}
	return err
}

func (repo *BaseRepository) AddCommand(channel string, commandToAdd *models.Command) error {
	repo.mutex.Lock()
	defer repo.mutex.Unlock()

	chanInfo, exists := repo.getChannelAndCommandExists(channel, commandToAdd)
	if exists {
		return CommandAlreadyExistsError
	}
	chanInfo.Commands[commandToAdd.Name] = *commandToAdd
	return repo.updateChannel(&chanInfo)
}

func (repo *BaseRepository) EditCommand(channel string, commandToEdit *models.Command) error {
	repo.mutex.Lock()
	defer repo.mutex.Unlock()

	chanInfo, exists := repo.getChannelAndCommandExists(channel, commandToEdit)
	if !exists {
		return CommandNameNotFoundError
	}
	chanInfo.Commands[commandToEdit.Name] = *commandToEdit
	return repo.updateChannel(&chanInfo)
}

func (repo *BaseRepository) DeleteCommand(channel string, commandToDelete *models.Command) error {
	repo.mutex.Lock()
	defer repo.mutex.Unlock()

	chanInfo, exists := repo.getChannelAndCommandExists(channel, commandToDelete)
	if !exists {
		return CommandNameNotFoundError
	}
	delete(chanInfo.Commands, commandToDelete.Name)
	return repo.updateChannel(&chanInfo)
}

// Returns channel info and if the command exists. Called by Add/Edit/Delete Command methods.
// This function assumes that mutex lock IS obtained by its caller.
func (repo *BaseRepository) getChannelAndCommandExists(
	channel string, command *models.Command) (models.Channel, bool) {

	chanInfo := repo.channelMap[channel]
	if chanInfo == nil {
		return models.Channel{}, false
	}
	_, exists := chanInfo.Commands[command.Name]
	return *chanInfo, exists // Make a copy not to make changes to the original
}

func (repo *BaseRepository) ListCommands(channel string) ([]*models.Command, error) {
	repo.mutex.RLock()
	defer repo.mutex.RUnlock()

	commandMap := repo.channelMap[channel].Commands
	commands := []*models.Command{}

	// Note: When iterating map, Go re-uses the same memory address for the values of the map.
	// If &v is appended to commands directly without copying, the slice will have pointers
	// to the same memory address, and therefore the same value.
	for _, v := range commandMap {
		copied := v // This copy is required
		commands = append(commands, &copied)
		// fmt.Println("copied: ", copied)
	}
	return commands, nil
}
