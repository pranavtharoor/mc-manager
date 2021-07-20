package config

import (
	"strings"

	"github.com/spf13/viper"
)

type Configuration struct {
	Bot BotConfiguration `mapstructure:"bot"`
}

type BotConfiguration struct {
	Token  string `mapstructure:"token"`
	Prefix string `mapstructure:"prefix"`
}

func setDefaults() {
	viper.SetDefault("bot.token", "")
	viper.SetDefault("bot.prefix", "!")
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
