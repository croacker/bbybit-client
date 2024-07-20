package config

import (
	"log"

	"github.com/spf13/viper"
)

type AppConfig struct {
	Symbols []string `mapstructure:"symbols"`
}

const CONFIG_FILE = "config.yml"
const CONFIG_TYPE = "yaml"

var appConfig *AppConfig

func GetAppConfig() *AppConfig {
	return appConfig
}

func LoadConfig() *AppConfig {
	log.Println("load configuration...")

	v := viper.New()
	v.SetConfigFile(CONFIG_FILE)
	v.SetConfigType(CONFIG_TYPE)

	if err := v.ReadInConfig(); err != nil {
		log.Fatal("error read config:", err)
	}

	appConfig = &AppConfig{}
	if err := v.Unmarshal(appConfig); err != nil {
		log.Fatal("error unmarshal config:", err)
	}

	log.Print("configuration:", *appConfig)

	return appConfig
}
