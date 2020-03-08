package manager

import (
	"sync"

	"github.com/c-rainbow/simplechatbot/client"
	"github.com/c-rainbow/simplechatbot/repository"
	twitch_irc "github.com/gempir/go-twitch-irc"
)

// ActivePluginManagerT interface for chat listener plugin manager
type ActivePluginManagerT interface {
	ProcessChat(channel string, sender *twitch_irc.User, message *twitch_irc.PrivateMessage)
	Close()
}

/*
ActivePluginManager manager for active plugins.

An active plugin is capable of responding in Twitch chat if necessary. It may also perform other actions in the chat,
like banning users from their chat history.

(1) Chats in the same Twitch channel are processed sequentially, in FIFO manner.
(2) For each chat, listener plugins run in order that they are registered in the manager.
(3) All listener plugins run first, then the command plugin manager runs.

The order of registration is important because a common usage of listener plugin is to punish
(ban, timeout, warn, etc) inappropriate chats. For example, a long emote-only chat can be punished for
(1) too long, or (2) spamming emotes. If one plugin decides to timeout the user, then other plugins should
not process the chat anymore.
*/
type ActivePluginManager struct {
	chanMapMutex          sync.RWMutex
	listenerPluginManager ChatListenerPluginManagerT
	commandPluginManager  ChatCommandPluginManagerT
	repo                  repository.SingleBotRepositoryT
	chanMap               map[string]chan WorkItem
}

var _ ActivePluginManagerT = (*ActivePluginManager)(nil)

// NewActivePluginManager creates a new ActivePluginManager
func NewActivePluginManager(
	repo repository.SingleBotRepositoryT, listenerPluginManager ChatListenerPluginManagerT,
	commandPluginManager ChatCommandPluginManagerT) ActivePluginManagerT {
	manager := ActivePluginManager{
		listenerPluginManager: listenerPluginManager, commandPluginManager: commandPluginManager,
		repo: repo, chanMap: make(map[string]chan WorkItem)}

	manager.initializeChanMap()

	return &manager
}

// DefaultActivePluginManager creates a new ActivePluginManager with IRC client and repository object
func DefaultActivePluginManager(
	ircClient client.TwitchClientT, repo repository.SingleBotRepositoryT) ActivePluginManagerT {
	listenerPluginManager := NewChatListenerPluginManager(ircClient, repo)
	commandPluginManager := NewChatCommandPluginManager(ircClient, repo)
	return NewActivePluginManager(repo, listenerPluginManager, commandPluginManager)
}

func (manager *ActivePluginManager) initializeChanMap() {
	manager.chanMapMutex.Lock()
	defer manager.chanMapMutex.Unlock()

	channels := manager.repo.GetAllChannels()
	for _, channel := range channels {
		jobChannel := make(chan WorkItem, jobChanDefaultCapacity)
		manager.chanMap[channel.Username] = jobChannel

		// The goroutine will automatically close when jobChannel is closed
		go func() {
			for workItem := range jobChannel {
				// First run all listeners, then run command-triggered plugins
				toContinue := manager.listenerPluginManager.ProcessChat(
					workItem.channel, workItem.sender, workItem.message)
				if toContinue {
					manager.commandPluginManager.ProcessChat(
						workItem.channel, workItem.sender, workItem.message)
				}
			}
		}()
	}
}

// ProcessChat parses command form chat message and calls the right plugin if necessary.
func (manager *ActivePluginManager) ProcessChat(
	channel string, sender *twitch_irc.User, message *twitch_irc.PrivateMessage) {

	manager.chanMapMutex.RLock()
	defer manager.chanMapMutex.RUnlock()
	manager.chanMap[channel] <- WorkItem{channel: channel, sender: sender, message: message}
}

// Close closes all job channels
func (manager *ActivePluginManager) Close() {
	// Close all channels
	manager.chanMapMutex.Lock()
	defer manager.chanMapMutex.Unlock()
	for _, channel := range manager.chanMap {
		close(channel)
	}
}
