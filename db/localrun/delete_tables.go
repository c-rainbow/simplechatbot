package localrun

import (
	"fmt"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/c-rainbow/simplechatbot/flags"
	"github.com/c-rainbow/simplechatbot/repository"
	"github.com/guregu/dynamo"
)

// Deletes all tables in local DynamoDB.
// WARNING: This command deletes all data and tables.
func DeleteAllTables() error {
	fmt.Println("endpoint: ", flags.DatabaseEndpoint)
	db := dynamo.New(session.New(), &aws.Config{
		Endpoint:   aws.String(flags.DatabaseEndpoint),
		Region:     aws.String(flags.DatabaseRegion),
		DisableSSL: aws.Bool(flags.DisableSSL),
	})

	var err error

	// Delete Bots table
	err = db.Table(repository.BotTableName).DeleteTable().Run()
	if err != nil {
		fmt.Println("Error while deleting Bots table. ", err.Error())
		return err
	}

	// Delete Channels table
	err = db.Table(repository.ChannelTableName).DeleteTable().Run()
	if err != nil {
		fmt.Println("Error while deleting Channels table. ", err.Error())
		return err
	}
	return nil
}
