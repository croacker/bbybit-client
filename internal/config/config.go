package config

import (
	"fmt"
	"log"

	"github.com/spf13/viper"
)

const CONFIG_FILE = "config.yml"
const CONFIG_TYPE = "yaml"

type AppConfig struct {
	BbClient BbClientConfig `mapstructure:"bb-client"`
	TgClient TgClientConfig `mapstructure:"tg-client"`
	Symbols  []string       `mapstructure:"symbols"`
}

type BbClientConfig struct {
	Url       string `mapstructure:"url"`
	ApiKey    string `mapstructure:"api-key"`
	ApiSecret string `mapstructure:"api-secret"`
	Interval  int32  `mapstructure:"interval-sec"`
}

type TgClientConfig struct {
	Token string `mapstructure:"token"`
}

var appConfig *AppConfig

func (c AppConfig) String() string {
	return fmt.Sprintf("{ bb-client: %v, tg-client: %v, symbols: %v }", c.BbClient, c.TgClient, c.Symbols)
}

func (b BbClientConfig) String() string {
	apiKey := ", api-key:*******"
	if len(b.ApiKey) == 0 {
		apiKey = ", api-key:EMPTY"
	}
	apiSecret := ", api-secret:*******"
	if len(b.ApiSecret) == 0 {
		apiSecret = ", api-secret:EMPTY"
	}
	return fmt.Sprintf("{ url: %v, interval-sec: %v %s %s }", b.Url, b.Interval, apiKey, apiSecret)
}

func (t TgClientConfig) String() string {
	token := "token:*******"
	if len(t.Token) == 0 {
		token = "token:EMPTY"
	}
	return fmt.Sprintf("{ %s }", token)
}

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

	log.Print("configuration loaded:", *appConfig)

	return appConfig
}
