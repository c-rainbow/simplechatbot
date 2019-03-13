// Create tables in local DynamoDB
package localrun

import (
	"fmt"

	"github.com/c-rainbow/simplechatbot"
	"github.com/c-rainbow/simplechatbot/models"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/guregu/dynamo"
)

// Create all necessary tables in DB.
// The function assumes that tables with the same names do not exist.
func CreateAllTables() error {
	db := dynamo.New(session.New(), &aws.Config{
		Endpoint:   aws.String(simplechatbot.DatabaseEndpoint),
		Region:     aws.String(simplechatbot.DatabaseRegion),
		DisableSSL: aws.Bool(simplechatbot.DisableSSL),
	})

	// Creates Bots table.
	err := db.
		CreateTable(simplechatbot.BotTableName, &models.Bot{}).          // Create table
		Project("Username-index", dynamo.IncludeProjection, "Username"). // Create index
		Run()
	if err != nil {
		fmt.Println("Error while creating Bots table: " + err.Error())
		return err
	}

	// Creates Channels table.
	err = db.
		CreateTable(simplechatbot.ChannelTableName, &models.Channel{}).  // Create table
		Project("Username-index", dynamo.IncludeProjection, "Username"). // Create index
		Run()
	if err != nil {
		fmt.Println("Error while creating Channels table: " + err.Error())
		return err
	}

	fmt.Println("Successfully creates all tables")
	return nil
}
