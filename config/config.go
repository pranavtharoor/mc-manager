package config

import (
	"strings"

	"github.com/spf13/viper"
)

type Configuration struct {
	Bot BotConfiguration `mapstructure:"bot"`
}

type BotConfiguration struct {
	Token      string                  `mapstructure:"token"`
	Prefix     string                  `mapstructure:"prefix"`
	Server     ServerConfiguration     `mapstructure:"server"`
	EasterEggs EasterEggsConfiguration `mapstructure:"easterEggs"`
}

type ServerConfiguration struct {
	ResourceGroup string `mapstructure:"resourceGroup"`
	Name          string `mapstructure:"name"`
}

type EasterEggsConfiguration struct {
	ReplyEgg ReplyEggConfiguration `mapstructure:"replyEgg"`
}

type ReplyEggConfiguration struct {
	Enabled  bool   `mapstructure:"enabled"`
	LookFor  string `mapstructure:"lookFor"`
	SayStart string `mapstructure:"sayStart"`
	SayEnd   string `mapstructure:"sayEnd"`
	TagUser  bool   `mapstructure:"tagUser"`
	ReplyTo  string `mapstructure:"replyTo"`
}

func setDefaults() {
	viper.SetDefault("bot.token", "")
	viper.SetDefault("bot.prefix", "!")
	viper.SetDefault("bot.server.resourceGroup", "")
	viper.SetDefault("bot.server.name", "")
	viper.SetDefault("bot.easterEggs.replyEgg.enabled", false)
	viper.SetDefault("bot.easterEggs.replyEgg.lookFor", "")
	viper.SetDefault("bot.easterEggs.replyEgg.sayStart", "")
	viper.SetDefault("bot.easterEggs.replyEgg.sayEnd", "")
	viper.SetDefault("bot.easterEggs.replyEgg.tagUser", true)
	viper.SetDefault("bot.easterEggs.replyEgg.replyTo", "")
}

func Read() (Configuration, error) {
	var c Configuration

	viper.SetConfigName("config")
	viper.AddConfigPath(".")
	viper.SetEnvPrefix("mc")
	viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))
	viper.AutomaticEnv()

	err := viper.ReadInConfig()
	if _, ok := err.(viper.ConfigFileNotFoundError); err != nil && !ok {
		return c, err
	}

	setDefaults()

	if err := viper.Unmarshal(&c); err != nil {
		return c, err
	}

	return c, nil
}
