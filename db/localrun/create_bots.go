package localrun

import (
	"fmt"

	simplechatbot "github.com/c-rainbow/simplechatbot"
	"github.com/c-rainbow/simplechatbot/models"
)

// Add new bot to local DynamoDB
func AddNewBot() {
	repo := simplechatbot.NewBaseRepository()

	// Make sure that no bots exist first
	bots := repo.GetAllBots()
	fmt.Println("Number of existing bots: ", len(bots))

	// Add one bot fixture
	testBot := &models.Bot{
		TwitchID:   DefaultBotTwitchID,
		Username:   DefaultBotUsername,
		OauthToken: DefaultBotOauthToken,
	}
	repo.CreateNewBot(testBot)
	bots = repo.GetAllBots()
	fmt.Println("New number of existing bots: ", len(bots))
	fmt.Println("Bot info: ", bots[0])
}
