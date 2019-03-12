package simplechatbot

import "flag"

var DatabaseEndpoint = *flag.String("dynamodb-endpoint", "http://localhost:8000", "DynamoDB endpoint address")
var DatabaseRegion = *flag.String("dynamodb-region", "us-west-2", "Default Region for DynamoDB")
var DisableSSL = *flag.Bool("dynamodb-disable-ssl", true, "If true, disable SSL to connect to DynamoDB")
