package games

// 인벤의 주사위게임 비슷하게 만들어보기
// http://www.inven.co.kr/board/cq/4006/167193
// 주사위 게임에 대해 알려주신 트수분 감사합니다.

import (
	"math/rand"
	"strconv"

	"github.com/c-rainbow/simplechatbot/client"
	"github.com/c-rainbow/simplechatbot/models"
	"github.com/c-rainbow/simplechatbot/parser"
	chatplugins "github.com/c-rainbow/simplechatbot/plugins/chat"
	"github.com/c-rainbow/simplechatbot/plugins/chat/common"
	twitch_irc "github.com/gempir/go-twitch-irc"
)

const (
	// DicePluginType plugin type name to roll a dice
	DicePluginType = "DicePluginType"

	maxDiceNumber = 100
)

var (
	diceResponseTexts = []string{
		"@$(user) 열심히 던져보았으나 아쉽게 $(arg0) 나왔습니다.",
		"@$(user) 대충 던진 주사위가 $(arg0) 이라니!!",
		"@$(user) $(arg0) 정도면 운이 좋은 건가요?",
		"주사위를 던졌다고 @$(user)님 처럼 아무나 $(arg0) 나오는건 아니지!",
		"@$(user) $(arg0) 나왔음",
	}
)

// DicePlugin plugin that responds to user-defined chat commands.
type DicePlugin struct {
	ircClient client.TwitchClientT
}

var _ chatplugins.ChatCommandPluginT = (*DicePlugin)(nil)

// NewDicePlugin creates a new DicePlugin
func NewDicePlugin(ircClient client.TwitchClientT) chatplugins.ChatCommandPluginT {
	return &DicePlugin{ircClient: ircClient}
}

// GetPluginType gets plugin type
func (plugin *DicePlugin) GetPluginType() string {
	return DicePluginType
}

// ReactToChat reacts to chat
func (plugin *DicePlugin) ReactToChat(
	command *models.Command, channel string, sender *twitch_irc.User, message *twitch_irc.PrivateMessage) {

	responseText := ""
	err := common.ValidateBasicInputs(command, channel, DicePluginType, sender, message)

	if err == nil {
		newNum := rand.Intn(maxDiceNumber) + 1
		index := rand.Intn(len(diceResponseTexts))
		response := diceResponseTexts[index]
		parsedResponse := parser.ParseResponse(response)
		args := []string{strconv.Itoa(newNum)}
		responseText, err = parser.ConvertResponse(parsedResponse, channel, sender, message, args)
	}

	common.SendToChatClient(plugin.ircClient, channel, responseText)
	common.HandleError(err)

}
