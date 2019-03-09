package main

import (
	// "fmt"

	"fmt"
	"strconv"

	"github.com/c-rainbow/simplechatbot"
	examples "github.com/c-rainbow/simplechatbot/db/examples"
	models "github.com/c-rainbow/simplechatbot/models"
)

func main() {
	mainChannel()
}

func mainChannel() {
	addChannel()
	repo := simplechatbot.NewBaseRepository()
	chans := repo.GetAllChannels()

	fmt.Println("Total chans found ", len(chans))

	for _, channel := range chans {
		fmt.Println(channel.TwitchID)
		fmt.Println(channel.Username)
		fmt.Println("displayname: " + channel.DisplayName)
		fmt.Println("command len: " + string(len(channel.Commands)))
		for _, botId := range channel.BotIDs {
			fmt.Println("Bot id: " + strconv.Itoa((int)(botId)))
		}
		for _, command := range channel.Commands {
			fmt.Println("Printing command")
			fmt.Println("Commandname: " + command.Name)
			fmt.Println("Command channelID: " + strconv.Itoa((int)(command.ChannelID)))
			fmt.Println("Command enabled: " + strconv.FormatBool(command.Enabled))
		}

	}
}

func addChannel() {
	db := examples.NewDatabase()
	chanTable := db.Table("Channels")

	commands := []models.Command{
		models.Command{
			Name:      "Command name",
			ChannelID: 4322,
		},
	}
	for _, command := range commands {
		fmt.Println("com name: " + command.Name)
		fmt.Println("com channelid: " + strconv.Itoa((int)(command.ChannelID)))
	}

	testChannel := &models.Channel{
		TwitchID:    6666,
		Username:    "teeeeee",
		DisplayName: "tdddddddd",
		BotIDs:      []int64{100, 200, 222200},
		/*Commands: []models.Command{
			models.Command{
				Name:      "Command name",
				ChannelID: 4322,
			},
		},*/
	}

	// putItem := chanTable.Put(testChannel)
	//putItem.input()
	err := chanTable.Put(testChannel).Run()
	if err != nil {
		fmt.Println("err: " + err.Error())
	}
}

func mainBot() {
	// addBot()
	repo := simplechatbot.NewBaseRepository()
	bots := repo.GetAllBots()

	fmt.Println("Total bots found ", len(bots))

	for _, bot := range bots { // range doesn't work with itr?
		fmt.Println(bot.TwitchID)
		fmt.Println(bot.Username)
		fmt.Println(bot.OauthToken)
	}
}

func addBot() {
	db := examples.NewDatabase()
	botTable := db.Table("Bots")
	botTable.Put(&models.Bot{
		TwitchID:   2345,
		Username:   "botName55",
		OauthToken: "dontknow22",
	}).Run()
}

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
/*
func main() {
	pack.InitializeTables()
}
*/
/*
func main() {
	response, err := http.Get("http://localhost:8000/shell")
	if err != nil {
		fmt.Printf("%s", err)
		os.Exit(1)
	} else {
		defer response.Body.Close()
		contents, err := ioutil.ReadAll(response.Body)
		if err != nil {
			fmt.Printf("%s", err)
			os.Exit(1)
		}
		fmt.Printf("%s\n", string(contents))
	}
}

/*
func main1() {
	// simplechatbot.GoMain()

	// anow := time.Now()
	// time.Now().Millisecond
	// fmt.Println("Now: " + strconv.FormatInt(anow.UnixNano(), 16))
	botInfo := simplechatbot.GetDefaultBotInfo()
	chatBot := simplechatbot.NewDefaultChatBot(botInfo)
	defer chatBot.Disconnect()
	go chatBot.Connect()
	fmt.Println("Conencted in main.go")
	time.Sleep(10 * time.Second)
	// chatBot.Disconnect()
	// fmt.Println("Disconnected in main.go")
}
*/
