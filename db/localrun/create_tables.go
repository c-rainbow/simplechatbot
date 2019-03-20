// Create tables in local DynamoDB
package localrun

import (
	"fmt"

	"github.com/c-rainbow/simplechatbot/flags"
	"github.com/c-rainbow/simplechatbot/models"
	"github.com/c-rainbow/simplechatbot/repository"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/guregu/dynamo"
)

// Create all necessary tables in DB.
// The function assumes that tables with the same names do not exist.
func CreateAllTables() error {
	db := dynamo.New(session.New(), &aws.Config{
		Endpoint:   aws.String(flags.DatabaseEndpoint),
		Region:     aws.String(flags.DatabaseRegion),
		DisableSSL: aws.Bool(flags.DisableSSL),
	})

	var err error

	// Creates Bots table.
	/*err := db.
		CreateTable(simplechatbot.BotTableName, &models.Bot{}).          // Create table
		Project("Username-index", dynamo.IncludeProjection, "Username"). // Create index
		Run()
	if err != nil {
		fmt.Println("Error while creating Bots table: " + err.Error())
		return err
	}*/

	// Creates Channels table.
	err = db.
		CreateTable(repository.ChannelTableName, &models.Channel{}).     // Create table
		Project("Username-index", dynamo.IncludeProjection, "Username"). // Create index
		Run()
	if err != nil {
		fmt.Println("Error while creating Channels table: " + err.Error())
		return err
	}

	fmt.Println("Successfully creates all tables")
	return nil
}
