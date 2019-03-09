package simplechatbot

import "flag"

var flagMap = map[string]interface{}{
	"DynamoEndpoint":   *flag.String("dynamodb-endpoint", "", "DynamoDB endpoint address"),
	"DynamoRegion":     *flag.String("dynamodb-region", "us-west-2", "Default Region for DynamoDB"),
	"DynamoDisableSSL": *flag.Bool("dynamodb-disable-ssl", true, "If true, disable SSL to connect to DynamoDB"),
}

func GetFlagString(name string) string {
	val, ok := flagMap[name].(string)
	if !ok {
		return ""
	}
	return val
}

func GetFlagBool(name string) bool {
	val, ok := flagMap[name].(bool)
	if !ok {
		return false
	}
	return val
}
