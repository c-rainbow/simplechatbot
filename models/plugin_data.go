package models

import (
	"strconv"
)

// PluginData struct of data to be used by plugins
type PluginData struct {
	// Partition key. Composite of (PluginType, ChannelID, Key)
	// TODO: Find out how to configure composite partition key in the library
	PrimaryKey string `dynamo:"ID,hash"`
	// Plugin Type
	PluginType string `index:"PluginType-index,hash"`
	// Channel's integer Twitch ID. Don't use username which can change
	ChannelID int64 `index:"ChannelID-index,hash"`
	// Key. Plugin-dependent. Chat ID, Viewer ID, etc
	Key string
	// Value the actualy type is plugin-dependent
	Value interface{}
}

// NewPluginData creates a new PluginData object
func NewPluginData(pluginType string, channelID int64, key string, value interface{}) *PluginData {
	primaryKey := pluginType + "-" + strconv.FormatInt(channelID, 10) + "-" + key
	return &PluginData{
		PrimaryKey: primaryKey,
		PluginType: pluginType,
		ChannelID:  channelID,
		Value:      value,
	}
}
