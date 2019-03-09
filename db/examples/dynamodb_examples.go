package examples

import (
	"fmt"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	models "github.com/c-rainbow/simplechatbot/models"
	"github.com/guregu/dynamo"
)

const (
	// TestEndpoint is Endpoint address for test
	TestEndpoint = "http://localhost:8000"
	// TestRegion Region for test
	TestRegion = "us-west-2"
)

// NewDatabase creates a new DB handle.
func NewDatabase() *dynamo.DB {

	return dynamo.New(session.New(), &aws.Config{
		Endpoint:   aws.String(TestEndpoint),
		Region:     aws.String(TestRegion),
		DisableSSL: aws.Bool(true),
	})
}

// CreateTableExample Example of CreateTable
func CreateTableExample() {
	db := NewDatabase()
	err := db.
		CreateTable("Bots", models.Bot{}).                               // Create table
		Project("Username-index", dynamo.IncludeProjection, "Username"). // Create index
		Run()                                                            // No change to DB without this.

	if err != nil {
		fmt.Println("error: " + err.Error())
	}
}

// PutExample tests adding new item to table.
func PutExample() {
	db := NewDatabase()
	chanTable := db.Table("Channels")
	chanTable.Put(&models.Channel{
		TwitchID: 12345,
		Username: "botName3",
		BotIDs:   []int64{1, 34},
	}).Run()
}

// ScanExample example Scan
func ScanExample() {
	db := NewDatabase()
	chanTable := db.Table("Channels")
	itr := chanTable.Scan().Iter() // Get all items in the table.
	channel := &models.Channel{}

	for i := 0; i < 10; i++ { // range doesn't work with itr?
		hasNext := itr.Next(&channel)
		if !hasNext {
			break
		}
		fmt.Println("channel: " + channel.Username + "...")
	}
}
