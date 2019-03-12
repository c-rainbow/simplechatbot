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
	db := dynamo.New(session.New(), &aws.Config{
		Endpoint:   aws.String("http://localhost:8000"),
		Region:     aws.String("us-west-1"),
		DisableSSL: aws.Bool(true),
	})

	// Delete Bots table
	err := db.Table(simplechatbot.BotTableName).DeleteTable().Run()
	if err != nil {
		fmt.Println("Error while deleting Bots table. ", err.Error())
		return err
	}

	// Delete Channels table
	err = db.Table(simplechatbot.ChannelTableName).DeleteTable().Run()
	if err != nil {
		fmt.Println("Error while deleting Channels table. ", err.Error())
		return err
	}

	return nil
}
