// Create tables in local DynamoDB
package localrun

import (
	"fmt"
	"time"

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
	/*err = db.
		CreateTable(repository.ChannelTableName, &models.Channel{}).     // Create table
		Project("Username-index", dynamo.IncludeProjection, "Username"). // Create index
		Run()
	if err != nil {
		fmt.Println("Error while creating Channels table: " + err.Error())
		return err
	}*/

	// Creates PluginData table.
	/*err = db.
		CreateTable(repository.PluginDataTableName, &models.PluginData{}).   // Create table
		Project("PluginType-index", dynamo.IncludeProjection, "PluginType"). // Create index
		Project("ChannelID-index", dynamo.IncludeProjection, "ChannelID").   // Create index
		Run()
	if err != nil {
		fmt.Println("Error while creating PluginData table: " + err.Error())
		return err
	}*/

	//err = createTable(db)

	data1 := models.PluginData{
		PrimaryKey: "TestPluginType1-1234",
		PluginType: "TestPluginType1",
		ChannelID:  1234,
		Key:        "testKey",
	}

	data2 := models.PluginData{
		PrimaryKey: "TestPluginType1-1235",
		PluginType: "TestPluginType1",
		ChannelID:  1235,
		Key:        "testKey2",
		Value:      12345678,
	}

	data3 := models.PluginData{
		PrimaryKey: "TestPluginType2-1234",
		PluginType: "TestPluginType2",
		ChannelID:  1234,
		Key:        "testKey3",
		Value:      []string{"hello", "world"},
	}

	//time.Sleep(3 * time.Second)

	pdTable := db.Table(repository.PluginDataTableName)
	err = pdTable.Put(data1).Run()
	if err != nil {
		fmt.Println("Error1: " + err.Error())
	}
	err = pdTable.Put(data2).Run()
	if err != nil {
		fmt.Println("Error2: " + err.Error())
	}
	err = pdTable.Put(data3).Run()
	if err != nil {
		fmt.Println("Error3: " + err.Error())
	}

	fmt.Println("Successfully creates all tables")
	return nil
}

func createTable(db *dynamo.DB) error {
	err := db.
		CreateTable(repository.PluginDataTableName, &models.PluginData{}).   // Create table
		Project("PluginType-index", dynamo.IncludeProjection, "PluginType"). // Create index
		Project("ChannelID-index", dynamo.IncludeProjection, "ChannelID").   // Create index
		Run()
	if err != nil {
		fmt.Println("Error while creating PluginData table: " + err.Error())
		return err
	}
	time.Sleep(5 * time.Second)
	return nil
}
