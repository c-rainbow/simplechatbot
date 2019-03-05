package main

import (
	"fmt"
	"strconv"
	"time"
)

/*
func getIRCClient() *Client {
	dataManager := &PrototypeDataManager{}
	botInfo := dataManager.getBotInfo()
	userInfo := dataManager.getUserInfo()
	client := NewClient(botInfo.name, botInfo.oauthToken)

	fmt.Println("ircc: " + string(ircc.Hello))

	client.Handlers.Add(girc.CONNECTED, func(c *girc.Client, e girc.Event) {
		// c.Cmd.SendRaw("CAP REQ :twitch.tv/tags")
		c.Cmd.SendRaw("CAP REQ :twitch.tv/tags twitch.tv/commands")

		c.Cmd.Join("#" + userInfo.username)

		// c.Cmd.Message("#c_rainbow", "hello h안녕하fgfhf세요ello")
	})

	client.Handlers.Add(girc.PING, func(c *girc.Client, e girc.Event) {
		fmt.Println("PING RECEIVED!!! " + e.Trailing)
	})

	client.Handlers.Add(girc.PRIVMSG, handlePrivMsg)

	return client
}

func main() {

	// addCommand("c_rainbow", "!welcome", "Welcome..", 5)

	client := getIRCClient()

	if err := client.Connect(); err != nil {
		log.Fatal(err)
	}

	fmt.Println(client.String())
}
*/

func main() {
	// simplechatbot.GoMain()

	anow := time.Now()
	// time.Now().Millisecond
	fmt.Println("Now: " + strconv.FormatInt(anow.UnixNano(), 16))
}
