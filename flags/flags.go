package flags

import (
	"flag"
)

var DatabaseEndpoint string
var DatabaseRegion string
var DisableSSL bool
var ResetChannels bool
var InstallationConfigFile string

func ParseAllFlags() {
	// In go, flags have to be defined before flag.Parse(), and their actual values can be accessed by pointers
	// only after flag.Parse() is called.
	endpointPtr := flag.String("dynamodb-endpoint", "https://dynamodb.us-west-2.amazonaws.com", "DynamoDB endpoint address")
	regionPtr := flag.String("dynamodb-region", "us-west-2", "Default Region for DynamoDB")
	sslPtr := flag.Bool("dynamodb-disable-ssl", false, "If true, disable SSL to connect to DynamoDB")
	resetPtr := flag.Bool("reset-channels", false, "If true, reset Channels table")
	configPtr := flag.String("install-config-file", "", "Path to config file to install Simplechatbot")

	flag.Parse()

	DatabaseEndpoint = *endpointPtr
	DatabaseRegion = *regionPtr
	DisableSSL = *sslPtr
	ResetChannels = *resetPtr
	InstallationConfigFile = *configPtr
}
