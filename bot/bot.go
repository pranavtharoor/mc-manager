package bot

import (
	"strings"

	"github.com/bwmarrin/discordgo"
	"github.com/pranavtharoor/mc-manager/config"
)

var botID string
var botConfig config.BotConfiguration

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

	bot.AddHandler(messageHandler)

	bot.Open()

	return bot.UpdateListeningStatus("'" + botConfig.Prefix + "'")
}

func messageHandler(s *discordgo.Session, m *discordgo.MessageCreate) {
	if m.Author.ID == botID || !strings.HasPrefix(m.Content, botConfig.Prefix) {
		return
	}

	msg := strings.TrimPrefix(m.Content, botConfig.Prefix)
	send := func(msg string) {
		s.ChannelMessageSend(m.ChannelID, msg)
	}
	words := strings.Fields(msg)
	maxCmdLength := 2
	for i := len(words); i < maxCmdLength; i++ {
		words = append(words, "")
	}

	switch words[0] {
	case "server":
		switch words[1] {
		case "start":
			send(serverStart(botConfig.Server))
		case "stop":
			send(serverStop(botConfig.Server))
		default:
			send(help())
		}
	case "azure":
		switch words[1] {
		case "login":
			azureLogin(send)
		case "logout":
			send(azureLogout())
		case "account":
			send(azureAccount())
		default:
			send(help())
		}
	default:
		send(help())
	}
}
