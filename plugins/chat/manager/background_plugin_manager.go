package manager

import (
	"sync"

	chatplugins "github.com/c-rainbow/simplechatbot/plugins/chat"
	"github.com/c-rainbow/simplechatbot/repository"
	twitch_irc "github.com/gempir/go-twitch-irc"
)

// BackgroundPluginManagerT interface for chat listener plugin manager
type BackgroundPluginManagerT interface {
	RegisterPlugin(plugin chatplugins.ChatListenerPluginT)
	ProcessChat(channel string, sender *twitch_irc.User, message *twitch_irc.PrivateMessage)
	Close()
}

/*
BackgroundPluginManager background plugin manager. All background plugins are registered here.

Background plugins listen to all chats in the channel, and are independent of one another. Action by
one background plugin does not affect another. Also, there is no guarantee that chats will be processed
in order of time.

Examples of background plugins are detecting username changes or updating last chat time.
*/
type BackgroundPluginManager struct {
	pluginsMutex sync.RWMutex
	repo         repository.SingleBotRepositoryT

	// Plugin type and its job chan should have the same index in slices.
	pluginTypes []string
	jobChans    []chan WorkItem
}

var _ BackgroundPluginManagerT = (*BackgroundPluginManager)(nil)

// NewBackgroundPluginManager creates a new BackgroundPluginManager
func NewBackgroundPluginManager(repo repository.SingleBotRepositoryT) BackgroundPluginManagerT {
	manager := BackgroundPluginManager{repo: repo, pluginTypes: []string{}, jobChans: []chan WorkItem{}}

	// TODO: creates plugins
	//manager.RegisterPlugin()

	return &manager
}

// RegisterPlugin registers a new plugin to the manager.
func (manager *BackgroundPluginManager) RegisterPlugin(plugin chatplugins.ChatListenerPluginT) {
	manager.pluginsMutex.Lock()
	defer manager.pluginsMutex.Unlock()

	// Check if the same type of plugin was already registered
	index := -1
	for i, pluginType := range manager.pluginTypes {
		if pluginType == plugin.GetPluginType() {
			index = i
			break
		}
	}

	var jobChan chan WorkItem
	if index == -1 {
		manager.pluginTypes = append(manager.pluginTypes, plugin.GetPluginType())
		jobChan = make(chan WorkItem, jobChanDefaultCapacity)
		manager.jobChans = append(manager.jobChans, jobChan)
	} else {
		jobChan = manager.jobChans[index]
	}

	// The goroutine will automatically close when jobChan is closed
	go func() {
		for workItem := range jobChan {
			plugin.ListenToChat(workItem.channel, workItem.sender, workItem.message)
		}
	}()
}

// ProcessChat parses command form chat message and calls the right plugin if necessary.
func (manager *BackgroundPluginManager) ProcessChat(
	channel string, sender *twitch_irc.User, message *twitch_irc.PrivateMessage) {

	manager.pluginsMutex.RLock()
	defer manager.pluginsMutex.RUnlock()

	for _, jobChan := range manager.jobChans {
		jobChan <- WorkItem{channel: channel, sender: sender, message: message}
	}
}

// Close closes all job channels
func (manager *BackgroundPluginManager) Close() {
	manager.pluginsMutex.Lock()
	defer manager.pluginsMutex.Unlock()

	for _, channel := range manager.jobChans {
		close(channel)
	}
}
