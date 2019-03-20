package chatplugins

import (
	"log"
	"strings"
	"sync"

	models "github.com/c-rainbow/simplechatbot/models"
	"github.com/c-rainbow/simplechatbot/repository"
	twitch_irc "github.com/gempir/go-twitch-irc"
)

/*
Chat command plugin manager. All chat command plugins are registered here,
and depends on the command, the manager distributes command struct to the right channel.

Each type of plugin has a channel and a worker pool with at least 1 worker, which keeps
polling from the channel until it closes.
*/

type ChatCommandPluginManagerT interface {
	RegisterPlugin(plugin ChatCommandPluginT, count int)
	ProcessChat(channel string, text string, sender *twitch_irc.User, message *twitch_irc.Message)
	Close()
}

type ChatCommandPluginManager struct {
	channelMapMutex sync.Mutex
	repo            repository.SingleBotRepositoryT
	chanMap         map[string]chan *WorkItem
}

var _ ChatCommandPluginManagerT = (*ChatCommandPluginManager)(nil)

type WorkItem struct {
	command *models.Command
	channel string
	text    string
	sender  *twitch_irc.User
	message *twitch_irc.Message
}

func (manager *ChatCommandPluginManager) RegisterPlugin(plugin ChatCommandPluginT, count int) {
	pluginType := plugin.GetPluginType()
	// Create job channel per plugin type
	manager.channelMapMutex.Lock()
	jobChannel, exists := manager.chanMap[pluginType]
	if !exists {
		jobChannel = make(chan *WorkItem)
		manager.chanMap[pluginType] = jobChannel
	}
	manager.channelMapMutex.Unlock()

	// The goroutine will automatically close when jobChannel is closed
	// TODO: Should I wait for the goroutine to close?
	go func() {
		for workItem := range jobChannel {
			plugin.Run(workItem.command.Name, workItem.channel, workItem.sender, workItem.message)
		}
	}()
}

func (manager *ChatCommandPluginManager) ProcessChat(
	channel string, text string, sender *twitch_irc.User, message *twitch_irc.Message) {
	// Get command name
	commandName := GetCommandName(message.Text)

	command := manager.repo.GetCommandByChannelAndName(channel, commandName)
	if command == nil { // Chat is not a bot command.
		return
	}

	pluginType := command.PluginType
	chanToAdd := manager.chanMap[pluginType]
	if chanToAdd != nil {
		chanToAdd <- &WorkItem{command: command, text: text, sender: sender, message: message}
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
