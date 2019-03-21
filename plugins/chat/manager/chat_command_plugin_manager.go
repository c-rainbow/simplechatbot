package manager

import (
	"errors"
	"log"
	"strings"
	"sync"

	"github.com/c-rainbow/simplechatbot/plugins/chat/games"

	"github.com/c-rainbow/simplechatbot/client"
	models "github.com/c-rainbow/simplechatbot/models"
	chatplugins "github.com/c-rainbow/simplechatbot/plugins/chat"
	"github.com/c-rainbow/simplechatbot/plugins/chat/commandplugins"
	"github.com/c-rainbow/simplechatbot/repository"
	twitch_irc "github.com/gempir/go-twitch-irc"
)

const (
	JobChanDefaultCapacity = 100
)

/*
Chat command plugin manager. All chat command plugins are registered here,
and depends on the command, the manager distributes command struct to the right channel.

Each type of plugin has a channel and a worker pool with at least 1 worker, which keeps
polling from the channel until it closes.
*/

type ChatCommandPluginManagerT interface {
	RegisterPlugin(factory chatplugins.ChatCommandPluginFactoryT, count int) error
	ProcessChat(channel string, sender *twitch_irc.User, message *twitch_irc.Message)
	Close()
}

type ChatCommandPluginManager struct {
	channelMapMutex sync.Mutex
	repo            repository.SingleBotRepositoryT
	chanMap         map[string]chan WorkItem
}

var _ ChatCommandPluginManagerT = (*ChatCommandPluginManager)(nil)

type WorkItem struct {
	command *models.Command
	channel string
	sender  *twitch_irc.User
	message *twitch_irc.Message
}

func NewChatCommandPluginManager(
	ircClient client.TwitchClientT, repo repository.SingleBotRepositoryT) ChatCommandPluginManagerT {
	manager := ChatCommandPluginManager{
		channelMapMutex: sync.Mutex{}, repo: repo, chanMap: make(map[string]chan WorkItem)}

	manager.RegisterPlugin(commandplugins.NewAddCommandPluginFactory(ircClient, repo), 1)
	manager.RegisterPlugin(commandplugins.NewDeleteCommandPluginFactory(ircClient, repo), 1)
	manager.RegisterPlugin(commandplugins.NewEditCommandPluginFactory(ircClient, repo), 1)
	manager.RegisterPlugin(commandplugins.NewListCommandsPluginFactory(ircClient, repo), 1)
	manager.RegisterPlugin(commandplugins.NewCommandResponsePluginFactory(ircClient, repo), 1)
	manager.RegisterPlugin(games.NewNumberGuesserPluginFactory(ircClient, repo), 1)

	return &manager
}

func (manager *ChatCommandPluginManager) RegisterPlugin(
	factory chatplugins.ChatCommandPluginFactoryT, count int) error {
	// Create job channel per plugin type
	pluginType := factory.GetPluginType()
	manager.channelMapMutex.Lock()
	jobChannel, exists := manager.chanMap[pluginType]
	if exists {
		return errors.New("same type of plugin is already registered")
	} else {
		jobChannel = make(chan WorkItem, JobChanDefaultCapacity)
		manager.chanMap[pluginType] = jobChannel
	}
	manager.channelMapMutex.Unlock()

	// The goroutine will automatically close when jobChannel is closed
	for i := 0; i < count; i++ {
		plugin := factory.BuildNewPlugin()
		go func() {
			for workItem := range jobChannel {
				plugin.Run(workItem.command, workItem.channel, workItem.sender, workItem.message)
			}
		}()
	}
	return nil
}

func (manager *ChatCommandPluginManager) ProcessChat(
	channel string, sender *twitch_irc.User, message *twitch_irc.Message) {

	// Parse command name and get command, if any
	commandName := GetCommandName(message.Text)
	command := manager.repo.GetCommandByChannelAndName(channel, commandName)
	if command == nil { // Chat is not a bot command.
		return
	}

	pluginType := command.PluginType
	chanToAdd := manager.chanMap[pluginType]
	if chanToAdd != nil {
		chanToAdd <- WorkItem{command: command, channel: channel, sender: sender, message: message}
	} else {
		// chanToAdd shouldn't be nil. This is software bug.
		log.Println("Something is wrong. chanToAdd is nil")
	}
}

func (manager *ChatCommandPluginManager) Close() {
	// Close all channels
	manager.channelMapMutex.Lock()
	defer manager.channelMapMutex.Unlock()
	for _, channel := range manager.chanMap {
		close(channel)
	}
}

func GetCommandName(text string) string {
	// strings.Fields deals with heading/trailing/non-space whitespaces.
	fields := strings.Fields(text)
	if len(fields) == 0 { // This is unlikely, just checking for malicious input
		return ""
	}
	return fields[0]
}
