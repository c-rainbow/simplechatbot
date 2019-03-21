package games

import (
	"math/rand"
	"strconv"
	"strings"
	"sync"

	"github.com/c-rainbow/simplechatbot/client"
	"github.com/c-rainbow/simplechatbot/models"
	chatplugins "github.com/c-rainbow/simplechatbot/plugins/chat"
	"github.com/c-rainbow/simplechatbot/plugins/chat/commandplugins"
	"github.com/c-rainbow/simplechatbot/repository"
	twitch_irc "github.com/gempir/go-twitch-irc"
)

/*
Chat command plugin which guesses number
For example,
  "!guess start" starts a new game
  "!guess stop" stops the current game
  "!guess restart" restarts with a new number
  "!guess [number] guesses a number.

During game, bot responses with one of "higher than [number]", "smaller than [number]", "[number] is correct!"
*/

const (
	NumberGuesserPluginType = "NumberGuesserPluginType"
	MaxNumber               = 100
	StatusRunning           = 1
	StatusStopped           = 0
)

// Plugin that responds to user-defined chat commands.
type NumberGuesserPlugin struct {
	ircClient client.TwitchClientT
	repo      repository.SingleBotRepositoryT
	mutex     sync.Mutex
	status    int
	answer    int
}

var _ chatplugins.ChatCommandPluginT = (*NumberGuesserPlugin)(nil)

func NewNumberGuesserPlugin(
	ircClient client.TwitchClientT, repo repository.SingleBotRepositoryT) chatplugins.ChatCommandPluginT {
	return &NumberGuesserPlugin{ircClient: ircClient, repo: repo, mutex: sync.Mutex{}}
}

func (plugin *NumberGuesserPlugin) GetPluginType() string {
	return NumberGuesserPluginType
}

func (plugin *NumberGuesserPlugin) Run(
	commandName string, channel string, sender *twitch_irc.User, message *twitch_irc.Message) error {
	// Read-action-print loop
	command, err := commandplugins.CommonRead(plugin.repo, commandName, channel, NumberGuesserPluginType, sender, message)
	toSay, err := plugin.action(command, channel, sender, message, err)
	err = commandplugins.CommonOutput(plugin.ircClient, channel, toSay, err)
	if err != nil {
		return err
	}
	return nil
}

// In this function, returned error means internal runtime error, not user input error.
// For example, NoPermissionsError is not an error here. However, a connection error to DB
// should be returned as error in this function.
//
// Note that CommandNotFoundError is also treated as an error, because in usual workflow,
// this plugin is only called after chat message handler checks for command name.
func (plugin *NumberGuesserPlugin) action(
	command *models.Command, channel string, sender *twitch_irc.User, message *twitch_irc.Message,
	err error) (string, error) {

	fields := strings.Fields(message.Text)

	toSay := ""

	argument := fields[1]
	plugin.mutex.Lock()
	defer plugin.mutex.Unlock()

	switch argument {
	case "시작":
		if plugin.status == StatusStopped {
			plugin.status = StatusRunning
			plugin.answer = rand.Intn(MaxNumber) + 1
			toSay = "1부터 " + strconv.Itoa(MaxNumber) + " 사이 제가 생각하는 숫자를 맞춰보세요"
		}
	case "종료":
		if plugin.status == StatusRunning {
			plugin.status = StatusStopped
			toSay = "게임을 종료합니다"
		}
	default:
		if plugin.status == StatusRunning {
			num, err := strconv.Atoi(argument)
			// non-nil err means invalid input. Ignore.
			if err == nil && 1 <= num && num <= MaxNumber {
				numStr := strconv.Itoa(num)
				if num < plugin.answer {
					toSay = numStr + " 보다 큽니다"
				} else if num > plugin.answer {
					toSay = numStr + "보다 작습니다"
				} else {
					toSay = "@" + sender.DisplayName + " 님 정답! 정답은 " + numStr + ". 게임을 종료합니다"
					plugin.status = StatusStopped
				}
			}
		}
	}
	return toSay, nil
}

type NumberGuesserPluginFactory struct {
	ircClient client.TwitchClientT
	repo      repository.SingleBotRepositoryT
}

var _ chatplugins.ChatCommandPluginFactoryT = (*NumberGuesserPluginFactory)(nil)

func NewNumberGuesserPluginFactory(
	ircClient client.TwitchClientT, repo repository.SingleBotRepositoryT) chatplugins.ChatCommandPluginFactoryT {
	return &NumberGuesserPluginFactory{ircClient: ircClient, repo: repo}
}

func (plugin *NumberGuesserPluginFactory) GetPluginType() string {
	return NumberGuesserPluginType
}

func (plugin *NumberGuesserPluginFactory) BuildNewPlugin() chatplugins.ChatCommandPluginT {
	return NewNumberGuesserPlugin(plugin.ircClient, plugin.repo)
}
