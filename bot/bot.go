package bot

import (
	"fmt"
	"regexp"
	"strings"
	"time"

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
	bot.AddHandler(easterEggHandler)
	bot.AddHandler(conversationHandler)

	bot.Open()

	ticker := time.NewTicker(5 * time.Hour)
	go func() {
		for {
			<-ticker.C
			_ = azureAccount()
		}
	}()

	return bot.UpdateListeningStatus("'" + botConfig.Prefix + "'")
}

func messageHandler(s *discordgo.Session, m *discordgo.MessageCreate) {
	m.Content = strings.ToLower(m.Content)
	if m.Author.ID == botID || !strings.HasPrefix(m.Content, strings.ToLower(botConfig.Prefix)) {
		return
	}
	s.ChannelTyping(m.ChannelID)
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

func easterEggHandler(s *discordgo.Session, m *discordgo.MessageCreate) {
	if botConfig.EasterEggs.ReplyEgg.Enabled {
		lookFor := strings.TrimSpace(botConfig.EasterEggs.ReplyEgg.LookFor)
		msg := strings.ToLower(m.Content)
		sayStart := botConfig.EasterEggs.ReplyEgg.SayStart
		sayEnd := botConfig.EasterEggs.ReplyEgg.SayEnd
		replyTo := botConfig.EasterEggs.ReplyEgg.ReplyTo
		tagUser := botConfig.EasterEggs.ReplyEgg.TagUser
		if matched, err := regexp.MatchString(lookFor, msg); m.Author.ID != botID && err == nil && matched {
			s.ChannelTyping(m.ChannelID)
			if replyTo != "" && replyTo == m.Author.ID {
				reply := ""
				if tagUser {
					reply = fmt.Sprintf("%s<@%s>%s", sayStart, replyTo, sayEnd)
				} else {
					reply = fmt.Sprintf("%s%s", sayStart, sayEnd)
				}
				s.ChannelMessageSend(m.ChannelID, reply)
			} else if tagUser {
				reply := fmt.Sprintf("%s<@%s>%s", sayStart, m.Author.ID, sayEnd)
				s.ChannelMessageSend(m.ChannelID, reply)
			} else {
				reply := fmt.Sprintf("%s%s", sayStart, sayEnd)
				s.ChannelMessageSend(m.ChannelID, reply)
			}
		}
	}
}

func conversationHandler(s *discordgo.Session, m *discordgo.MessageCreate) {
	isMentioned := false
	for _, user := range m.Mentions {
		if user.Bot {
			isMentioned = true
			break
		}
	}
	if !isMentioned || m.Author.ID == botID {
		return
	}
	s.ChannelTyping(m.ChannelID)
	messages, err := s.ChannelMessages(m.ChannelID, botConfig.Conversation.ContextLength+1, "", "", "")
	if err != nil {
		s.ChannelMessageSend(m.ChannelID, err.Error())
	}
	if len(messages) > 0 {
		messages = messages[:len(messages)-1]
	}
	replacer := regexp.MustCompile(`<@\!.*?>`)
	var pastUserInputs []string
	var generatedResponses []string
	for i := len(messages) - 1; i >= 0; i-- {
		message := messages[i]
		isMentioned := false
		for _, user := range message.Mentions {
			if user.Bot {
				isMentioned = true
				break
			}
		}
		content := replacer.ReplaceAllString(message.Content, "")
		if message.Author.ID == botID {
			generatedResponses = append(generatedResponses, content)
		} else if isMentioned {
			pastUserInputs = append(pastUserInputs, content)
		}
	}
	text := replacer.ReplaceAllString(m.Content, "")
	s.ChannelMessageSend(m.ChannelID, conversation(botConfig.Conversation, text, pastUserInputs, generatedResponses, 0))
}
