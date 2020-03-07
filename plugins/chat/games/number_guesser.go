package games

import (
	"fmt"
	"math"
	"math/rand"
	"strconv"
	"strings"
	"sync"

	"github.com/c-rainbow/simplechatbot/client"
	"github.com/c-rainbow/simplechatbot/models"
	"github.com/c-rainbow/simplechatbot/parser"
	chatplugins "github.com/c-rainbow/simplechatbot/plugins/chat"
	"github.com/c-rainbow/simplechatbot/plugins/chat/common"
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
	// NumberGuesserPluginType plugin type name to guess numbers
	NumberGuesserPluginType = "NumberGuesserPluginType"

	maxNumber     = 100
	startCommand  = "시작"
	endCommand    = "종료"
	statusRunning = 1
	statusStopped = 0
)

var (
	// Below are bot's response messages in various situations
	messageUsageBeforeGame    = "채팅창에 '!숫자 시작' 이라고 입력하여 게임을 시작할 수 있습니다"
	messageUsageInGame        = "채팅창에 '!숫자 [숫자]' 라고 입력하여 1부터 $(arg1) 사이 제가 생각하는 숫자를 맞춰보세요"
	messageGameStarted        = "게임이 시작되었습니다. " + messageUsageInGame
	messageGameAlreadyStarted = "이미 " + messageGameStarted
	messageLowerThanThat      = "$(arg0) 보다 작습니다"
	messageHigherThanThat     = "$(arg0) 보다 큽니다"
	messageCorrect            = "$(user)님 정답! 정답은 $(arg0)이었습니다. 게임을 종료합니다"
	messageGameEnded          = "게임을 종료합니다."
	messageGameAlreadyEnded   = "게임이 진행중이 아닙니다. " + messageUsageBeforeGame
)

// NumberGuesserPlugin a game plugin that
type NumberGuesserPlugin struct {
	ircClient  client.TwitchClientT
	repo       repository.SingleBotRepositoryT
	mutex      sync.Mutex
	status     int
	answer     int
	currentMax int
}

var _ chatplugins.ChatCommandPluginT = (*NumberGuesserPlugin)(nil)

// NewNumberGuesserPlugin creates a new NumberGuesserPlugin
func NewNumberGuesserPlugin(
	ircClient client.TwitchClientT, repo repository.SingleBotRepositoryT) chatplugins.ChatCommandPluginT {
	return &NumberGuesserPlugin{ircClient: ircClient, repo: repo, mutex: sync.Mutex{}}
}

// GetPluginType gets plugin type
func (plugin *NumberGuesserPlugin) GetPluginType() string {
	return NumberGuesserPluginType
}

// ReactToChat reacts to chat
func (plugin *NumberGuesserPlugin) ReactToChat(
	command *models.Command, channel string, sender *twitch_irc.User, message *twitch_irc.PrivateMessage) {
	responseText := ""
	err := common.ValidateBasicInputs(command, channel, NumberGuesserPluginType, sender, message)
	if err == nil {
		args := plugin.parseArguments(message.Message)

		response := plugin.processArgument(args, sender)
		// Parse the response message from above
		parsedResponse := parser.ParseResponse(response)

		num, err := strconv.ParseFloat(args[0], 64)
		responseArgs := []string{args[0], strconv.Itoa(plugin.currentMax)}
		if err == nil {
			num = math.Round(num)
			converted := fmt.Sprintf("%.0f", num)
			responseArgs[0] = converted
		}

		responseText, err = parser.ConvertResponse(parsedResponse, channel, sender, message, responseArgs)
	}

	common.SendToChatClient(plugin.ircClient, channel, responseText)
	common.HandleError(err)

}

func (plugin *NumberGuesserPlugin) parseArguments(text string) []string {
	fields := strings.Fields(text)

	switch len(fields) {
	case 0, 1:
		return []string{"", ""}
	case 2:
		return []string{fields[1], ""}
	default:
		return fields[1:]
	}
}

func (plugin *NumberGuesserPlugin) processArgument(args []string, sender *twitch_irc.User) string {
	plugin.mutex.Lock()
	defer plugin.mutex.Unlock()

	if plugin.status == statusRunning {
		return plugin.processInGameCommands(args)
	}
	return plugin.processNotInGameCommands(args)
}

// What to do while game is being played
func (plugin *NumberGuesserPlugin) processInGameCommands(args []string) string {
	mainArg := args[0]
	toSay := ""

	switch mainArg {

	// Duplicate start command. Show usage and ignore
	case startCommand:
		toSay = messageGameAlreadyStarted

	// End command. End the game.
	case endCommand:
		toSay = messageGameEnded
		plugin.status = statusStopped

	// Hopefully the chatter entered a valid number
	default:
		// Parse the number
		//num, err := strconv.Atoi(mainArg)
		num, err := strconv.ParseFloat(mainArg, 64)
		currentMaxFloat := float64(plugin.currentMax)
		answerFloat := float64(plugin.answer)
		if err == nil {
			num = math.Round(num)
		}
		if err == nil && 1 <= num && num <= currentMaxFloat {
			if num < answerFloat { // answer higher than number
				toSay = messageHigherThanThat
			} else if num > answerFloat { // answer lower than number
				toSay = messageLowerThanThat
			} else { // correct!
				toSay = messageCorrect
				plugin.status = statusStopped
			}
		} else {
			// Invalid number. Just show usage
			toSay = messageUsageInGame
		}
	}
	return toSay
}

// What to do when game is not being played
func (plugin *NumberGuesserPlugin) processNotInGameCommands(args []string) string {
	mainArg := args[0]
	toSay := ""

	switch mainArg {

	// Start the game.
	case startCommand:
		toSay = messageGameStarted
		// Check if user provided custom max number
		maxNum := maxNumber
		maxNum, err := strconv.Atoi(args[1])
		if err != nil {
			maxNum = maxNumber
		}
		// Set values for new game
		plugin.currentMax = maxNum
		plugin.answer = rand.Intn(maxNum) + 1
		plugin.status = statusRunning

	// While not in game, anything other than start command is invalid.
	default:
		toSay = messageGameAlreadyEnded
	}
	return toSay
}
