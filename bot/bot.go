package bot

import (
	"strings"

	"github.com/bwmarrin/discordgo"
	"github.com/pranavtharoor/mc-manager/config"
)

var botID string
var botPrefix string

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
	botPrefix = c.Prefix

	bot.AddHandler(messageHandler)

	bot.Open()

	return nil
}

func messageHandler(s *discordgo.Session, m *discordgo.MessageCreate) {
	if m.Author.ID == botID || !strings.HasPrefix(m.Content, botPrefix) {
		return
	}

	msg := strings.TrimPrefix(m.Content, botPrefix)
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
			send("Unimplemented")
		case "stop":
			send("Unimplemented")
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
