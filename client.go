package simplechatbot

import (
	"fmt"
	"log"

	"github.com/lrstanley/girc"
)

func getIRCClient() *Client {

	dataManager := &PrototypeDataManager{}
	botInfo := dataManager.getBotInfo()

	fmt.Println("bot info: " + botInfo.Name)

	userInfo := dataManager.getUserInfo()
	client := NewClient(botInfo.Name, botInfo.OauthToken)
	messageHandler := MessageHandler{}

	client.Handlers.Add(girc.CONNECTED, func(c *girc.Client, e girc.Event) {
		// This is needed to include tags in chat raw messages
		c.Cmd.SendRaw("CAP REQ :twitch.tv/tags twitch.tv/commands")
		c.Cmd.Join("#" + userInfo.Username)
	})

	client.Handlers.Add(girc.PRIVMSG, messageHandler.handlePrivmsg)

	return client
}

func GoMain() {

	// addCommand("c_rainbow", "!welcome", "Welcome..", 5)

	client := getIRCClient()

	if err := client.Connect(); err != nil {
		log.Fatal(err)
	}

	fmt.Println(client.String())
}
