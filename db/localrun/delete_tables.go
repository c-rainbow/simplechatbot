package localrun

import (
	"fmt"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/c-rainbow/simplechatbot"
	"github.com/guregu/dynamo"
)

// Deletes all tables in local DynamoDB.
// WARNING: This command deletes all data and tables.
func DeleteAllTables() error {
	fmt.Println("endpoint: ", simplechatbot.DatabaseEndpoint)
	db := dynamo.New(session.New(), &aws.Config{
		Endpoint:   aws.String(simplechatbot.DatabaseEndpoint),
		Region:     aws.String(simplechatbot.DatabaseRegion),
		DisableSSL: aws.Bool(simplechatbot.DisableSSL),
	})

	var err error

	// Delete Bots table
	/*err := db.Table(simplechatbot.BotTableName).DeleteTable().Run()
	if err != nil {
		fmt.Println("Error while deleting Bots table. ", err.Error())
		return err
	}*/

	// Delete Channels table
	err = db.Table(simplechatbot.ChannelTableName).DeleteTable().Run()
	if err != nil {
		fmt.Println("Error while deleting Channels table. ", err.Error())
		return err
	}

	return nil
}
