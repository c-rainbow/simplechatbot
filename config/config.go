package config

import (
	"fmt"
	"log"

	"github.com/c-rainbow/simplechatbot/flags"
	"github.com/go-ini/ini"
)

/*type Config struct {
	cfg *ini.File
}*/

func NewConfig() (*ini.File, error) {
	configFile := flags.InstallationConfigFile
	fmt.Println("Config file: ", configFile)
	cfg, err := ini.Load(flags.InstallationConfigFile)
	if err != nil {
		log.Println("Error while loading config file.", err.Error())
		return nil, err
	}
	return cfg, nil
}
