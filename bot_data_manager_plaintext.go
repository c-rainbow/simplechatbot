package simplechatbot

import (
	"fmt"
	"io/ioutil"
)

type PlainTextBotDataManager struct {
	botInfo   *Bot
	userInfos []*User
}

var _ BotDataManager = &PlainTextBotDataManager{}

func GetDefaultBotInfo() *Bot {
	botUsername, _ := ioutil.ReadFile("c:/Data/bot_username.txt")
	oauthToken, _ := ioutil.ReadFile("c:/Data/bot_oauth.txt")

	fmt.Println("bot username: " + string(botUsername))
	fmt.Println("oauth: " + string(oauthToken))

	return &Bot{
		Name:       string(botUsername),
		OauthToken: string(oauthToken),
	}
}

func (manager *PlainTextBotDataManager) GetBotInfo() *Bot {
	if manager.botInfo != nil {
		return manager.botInfo
	}

	manager.botInfo = GetDefaultBotInfo()

	return manager.botInfo
}

func (manager *PlainTextBotDataManager) GetUserInfos() []*User {

	if manager.userInfos != nil {
		return manager.userInfos
	}

	streamerUsername, _ := ioutil.ReadFile("c:/Data/streamer_username.txt")
	fmt.Println("streamer username: " + string(streamerUsername))
	manager.userInfos = []*User{
		&User{
			Username: string(streamerUsername),
		},
	}

	return manager.userInfos
}

func (manager *PlainTextBotDataManager) GetCommandsPerUser(user *User) []*Command {
	streamerUsername, _ := ioutil.ReadFile("c:/Data/streamer_username.txt")
	if user.Username == string(streamerUsername) {
		return []*Command{
			&Command{
				Name:     "hello",
				Response: "hello back",
			},
		}
	}

	return []*Command{}
}

func (manager *PlainTextBotDataManager) GetAllUserCommands() []*UserCommand {
	streamerUsername, _ := ioutil.ReadFile("c:/Data/streamer_username.txt")
	return []*UserCommand{
		&UserCommand{
			User: &User{
				Username: string(streamerUsername),
			},
			Commands: []*Command{
				&Command{
					Name:     "hello",
					Response: "hello back",
				},
			},
		},
	}
}

func (manager *PlainTextBotDataManager) AddCommand(user *User, command *Command) {

}

func (manager *PlainTextBotDataManager) EditCommand(user *User, command *Command) {

}

func (manager *PlainTextBotDataManager) DeleteCommand(user *User, command *Command) {

}
