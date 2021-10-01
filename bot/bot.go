package bot

import (
	"github.com/bwmarrin/discordgo"
	"github.com/pranavtharoor/mc-manager/config"
	"github.com/pranavtharoor/mc-manager/music"
)

var botID string
var botConfig config.BotConfiguration
var stereos map[string]*music.Stereo

func Start(c config.BotConfiguration) error {
	bot, err := discordgo.New("Bot " + c.Token)
	if err != nil {
		return err
	}

	u, err := bot.User("@me")
	if err != nil {
		return err
	}

	botID = u.ID
	botConfig = c
	stereos = map[string]*music.Stereo{}

	bot.AddHandler(helpHandler)
	bot.AddHandler(serverHandler)
	bot.AddHandler(easterEggHandler)
	bot.AddHandler(conversationHandler)
	bot.AddHandler(djHandler)

	bot.Open()

	return bot.UpdateListeningStatus("'" + botConfig.Prefix + "'")
}
