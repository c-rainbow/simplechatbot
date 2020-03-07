package manager

import (
	"log"
	"strings"
	"sync"

	"github.com/c-rainbow/simplechatbot/plugins/chat/games"
	"github.com/c-rainbow/simplechatbot/plugins/chat/selfban"

	"github.com/c-rainbow/simplechatbot/client"
	"github.com/c-rainbow/simplechatbot/models"
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
	RegisterPlugin(plugin chatplugins.ChatCommandPluginT)
	ProcessChat(channel string, sender *twitch_irc.User, message *twitch_irc.PrivateMessage)
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
	message *twitch_irc.PrivateMessage
}

func NewChatCommandPluginManager(
	ircClient client.TwitchClientT, repo repository.SingleBotRepositoryT) ChatCommandPluginManagerT {
	manager := ChatCommandPluginManager{
		channelMapMutex: sync.Mutex{}, repo: repo, chanMap: make(map[string]chan WorkItem)}

	manager.RegisterPlugin(commandplugins.NewAddCommandPlugin(ircClient, repo))
	manager.RegisterPlugin(commandplugins.NewDeleteCommandPlugin(ircClient, repo))
	manager.RegisterPlugin(commandplugins.NewEditCommandPlugin(ircClient, repo))
	manager.RegisterPlugin(commandplugins.NewListCommandsPlugin(ircClient, repo))
	manager.RegisterPlugin(commandplugins.NewCommandResponsePlugin(ircClient, repo))
	manager.RegisterPlugin(games.NewNumberGuesserPlugin(ircClient, repo))
	manager.RegisterPlugin(selfban.NewBanOnselfPlugin(ircClient))
	manager.RegisterPlugin(games.NewDicePlugin(ircClient))

	return &manager
}

// RegisterPlugin registers a new plugin to the manager.
// It is possible to register the same type of plugin multiple times, in which case there will be multiple workers
// for the same command.
func (manager *ChatCommandPluginManager) RegisterPlugin(plugin chatplugins.ChatCommandPluginT) {
	// Create job channel per plugin type
	pluginType := plugin.GetPluginType()
	manager.channelMapMutex.Lock()
	jobChannel, exists := manager.chanMap[pluginType]
	if !exists {
		jobChannel = make(chan WorkItem, JobChanDefaultCapacity)
		manager.chanMap[pluginType] = jobChannel
	}
	manager.channelMapMutex.Unlock()

	// The goroutine will automatically close when jobChannel is closed
	go func() {
		for workItem := range jobChannel {
			plugin.ReactToChat(workItem.command, workItem.channel, workItem.sender, workItem.message)
		}
	}()
}

// ProcessChat parses command form chat message and calls the right plugin if necessary.
func (manager *ChatCommandPluginManager) ProcessChat(
	channel string, sender *twitch_irc.User, message *twitch_irc.PrivateMessage) {

	// Parse command name and get command, if any
	commandName := getCommandName(message.Message)
	command := manager.repo.GetCommandByChannelAndName(channel, commandName)
	if command == nil { // Chat is not a bot command.
		return
	}

	pluginType := command.PluginType
	chanToAdd, exists := manager.chanMap[pluginType]
	if exists {
		chanToAdd <- WorkItem{command: command, channel: channel, sender: sender, message: message}
	} else {
		// chanToAdd shouldn't be nil. This is software bug.
		log.Println("Something is wrong. chanToAdd is nil")
	}
}

// Close closes all job channels
func (manager *ChatCommandPluginManager) Close() {
	// Close all channels
	manager.channelMapMutex.Lock()
	defer manager.channelMapMutex.Unlock()
	for _, channel := range manager.chanMap {
		close(channel)
	}
}

func getCommandName(text string) string {
	// strings.Fields deals with heading/trailing/non-space whitespaces.
	fields := strings.Fields(text)
	if len(fields) == 0 { // This is unlikely, just checking for malicious input
		return ""
	}
	return fields[0]
}
