package localrun

import (
	"fmt"

	"github.com/c-rainbow/simplechatbot/models"
	"github.com/c-rainbow/simplechatbot/repository"
)

// Add new bot to local DynamoDB
func AddNewBot() {
	repo := repository.NewBaseRepository()

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
