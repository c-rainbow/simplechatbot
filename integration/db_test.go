package integration

import (
	"fmt"
	"log"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/c-rainbow/simplechatbot/models"
	"github.com/c-rainbow/simplechatbot/repository"
	"github.com/guregu/dynamo"
)

// Gets local DynamoDB and cleans the entire data.
// WARNING: This command deletes all data in the local DynamoDB
func PreparaeCleanDb() error {
	db := dynamo.New(session.New(), &aws.Config{
		Endpoint:   aws.String("http://localhost:8000"),
		Region:     aws.String("us-west-1"),
		DisableSSL: aws.Bool(true),
	})
	tableNames, err := db.ListTables().All()

	fmt.Println("Existing table names:", tableNames)

	if err != nil {
		log.Fatal("Could not list tables in DB. " + err.Error())
		return err
	}

	if !inArray(tableNames, repository.BotTableName) {
		// Create Bots table if not exists
		createTable(db, repository.BotTableName, &models.Bot{})
	}

	if !inArray(tableNames, repository.ChannelTableName) {
		// Create Bots table if not exists
		createTable(db, repository.ChannelTableName, &models.Channel{})
	}
	return nil
}

func createTable(db *dynamo.DB, tableName string, model interface{}) error {
	fmt.Println("Table", tableName, "does not exist. Creating...")
	err := db.
		CreateTable(tableName, &model).                                  // Create table
		Project("Username-index", dynamo.IncludeProjection, "Username"). // Create index
		Run()                                                            // No change to DB without this.

	if err != nil {
		fmt.Println("error: " + err.Error())
		return err
	}
	return nil
}

func inArray(tableNames []string, tableName string) bool {
	for _, name := range tableNames {
		if name == tableName {
			return true
		}
	}
	return false
}

func TestPut(t *testing.T) {
	PreparaeCleanDb()
	repo := &repository.BaseRepository{}

	// Make sure that no bots exist first
	bots := repo.GetAllBots()
	assert.Empty(t, bots)

	// Add one bot
	testBot := &models.Bot{TwitchID: 1234, Username: "test_bot", OauthToken: "asdf1234"}
	repo.CreateNewBot(testBot)
	bots = repo.GetAllBots()
	assert.Equal(t, bots, []*models.Bot{testBot})
}

func TestSomething(t *testing.T) {

}
