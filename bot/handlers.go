package bot

import (
	"fmt"
	"regexp"
	"strconv"
	"strings"

	"github.com/bwmarrin/discordgo"
	"github.com/pranavtharoor/mc-manager/music"
)

func helpHandler(s *discordgo.Session, m *discordgo.MessageCreate) {
	topLevelCommands := []string{"server", "azure", "dj"}
	m.Content = strings.ToLower(m.Content)
	if m.Author.ID == botID || !strings.HasPrefix(m.Content, strings.ToLower(botConfig.Prefix)) {
		return
	}
	msg := strings.TrimPrefix(m.Content, botConfig.Prefix)
	send := func(msg string) {
		s.ChannelMessageSend(m.ChannelID, msg)
	}
	words := strings.Fields(msg)
	if len(words) > 0 {
		topLevelCommand := words[0]
		for _, command := range topLevelCommands {
			if command == topLevelCommand {
				return
			}
		}
	}
	s.ChannelTyping(m.ChannelID)
	send(help())
}

func serverHandler(s *discordgo.Session, m *discordgo.MessageCreate) {
	m.Content = strings.ToLower(m.Content)
	if m.Author.ID == botID || !strings.HasPrefix(m.Content, strings.ToLower(botConfig.Prefix)) {
		return
	}
	msg := strings.TrimPrefix(m.Content, botConfig.Prefix)
	send := func(msg string) {
		s.ChannelMessageSend(m.ChannelID, msg)
	}
	words := strings.Fields(msg)

	if len(words) > 0 {
		command := words[0]
		found := false
		for _, c := range []string{"server", "azure"} {
			if command == c {
				found = true
			}
		}
		if !found {
			return
		}
	}

	s.ChannelTyping(m.ChannelID)

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
		if user.Bot && user.ID == botID {
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
			if user.Bot && user.ID == botID {
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

func djHandler(s *discordgo.Session, m *discordgo.MessageCreate) {
	input := ""
	index := 0

	originalCase := m.Content
	m.Content = strings.ToLower(m.Content)
	if m.Author.ID == botID || !strings.HasPrefix(m.Content, strings.ToLower(botConfig.Prefix)) {
		return
	}
	msg := strings.TrimPrefix(m.Content, botConfig.Prefix)
	originalCase = originalCase[len(originalCase)-len(msg):]

	send := func(msg string) {
		if msg != "" {
			s.ChannelMessageSend(m.ChannelID, msg)
		}
	}
	words := strings.Fields(msg)
	originalCaseWords := strings.Fields(originalCase)

	if len(words) > 0 {
		command := words[0]
		found := false
		for _, c := range []string{"dj"} {
			if command == c {
				found = true
			}
		}
		if !found {
			return
		}
	}

	stereo, ok := stereos[m.GuildID]
	if !ok {
		getVoiceConnection := func() *discordgo.VoiceConnection {
			return s.VoiceConnections[m.GuildID]
		}
		stereo = music.NewStereo(getVoiceConnection)
		stereos[m.GuildID] = stereo
	}

	switch words[0] {
	case "dj":
		switch words[1] {
		case "play":
			fallthrough
		case "add":
			if len(words) < 3 {
				send("Give me a song")
				return
			}
			input = strings.Join(originalCaseWords[2:], " ")
		case "insert":
			fallthrough
		case "remove":
			if len(words) < 3 {
				send("Give me a index and song")
				return
			}
			var err error
			index, err = strconv.Atoi(originalCaseWords[2])
			if err != nil {
				send("Give me a index")
				return
			}
			index = index - 1
			input = strings.Join(originalCaseWords[2:], " ")
		}
	}

	switch words[0] {
	case "dj":
		switch words[1] {
		case "join":
			djJoin(s, m, send)
		case "leave":
			djLeave(send, stereo)
		case "play":
			s.ChannelTyping(m.ChannelID)
			djPlay(botConfig.Dj, stereo, input, send)
		case "add":
			s.ChannelTyping(m.ChannelID)
			djAdd(botConfig.Dj, stereo, input, send, len(stereo.GetQueue()))
		case "insert":
			s.ChannelTyping(m.ChannelID)
			djAdd(botConfig.Dj, stereo, input, send, index)
		case "remove":
			s.ChannelTyping(m.ChannelID)
			djRemove(stereo, send, index)
		case "list":
			s.ChannelTyping(m.ChannelID)
			send(djList(stereo))
		case "skip":
			stereo.Next()
		case "clear":
			s.ChannelTyping(m.ChannelID)
			stereo.ClearQueue()
			send("Queue Cleared")
		default:
			send(help())
		}
	}
}
