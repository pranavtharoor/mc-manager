package bot

import (
	"fmt"

	"github.com/bwmarrin/discordgo"
	"github.com/pranavtharoor/mc-manager/config"
	"github.com/pranavtharoor/mc-manager/music"
)

func djJoin(s *discordgo.Session, m *discordgo.MessageCreate, send func(string)) {
	var voiceChannelId string
	for _, g := range s.State.Guilds {
		for _, v := range g.VoiceStates {
			// TODO: check if user is in the voice channel
			voiceChannelId = v.ChannelID
		}
	}
	if voiceChannelId == "" {
		send("You should join a voice channel before asking me to")
		return
	}
	_, err := s.ChannelVoiceJoin(m.GuildID, voiceChannelId, false, false)

	if err != nil {
		send("Sorry, I couldn't join")
	}
}

func djLeave(send func(string), stereo *music.Stereo) {
	if !stereo.IsConnected() {
		send("I'm not in any voice channels")
		return
	}
	stereo.Disconnect()
}

func djPlay(c config.DjConfiguration, stereo *music.Stereo, searchOrUrl string, send func(string)) {
	if !stereo.IsConnected() {
		send("I'm not in any voice channels")
		return
	}

	song, err := getSongFromSearchOrUrl(c, searchOrUrl)
	if err != nil {
		send(err.Error())
		return
	}

	send("Playing " + song.Title)
	err = stereo.StartStreaming(song)
	if err != nil {
		send(err.Error())
	}
}

func djList(stereo *music.Stereo) string {
	songs := stereo.GetQueue()
	output := "Queue:\n"
	for i, song := range songs {
		output += fmt.Sprintf("\n%d. %s", i+1, song.Title)
	}
	return output
}

func djAdd(c config.DjConfiguration, stereo *music.Stereo, searchOrUrl string, send func(string), index int) {
	song, err := getSongFromSearchOrUrl(c, searchOrUrl)
	if err != nil {
		send("Un-oh, " + err.Error())
		return
	}

	send("Adding " + song.Title)
	stereo.AddToQueue(song, index)
}

func djRemove(stereo *music.Stereo, send func(string), index int) {
	if err := stereo.RemoveFromQueue(index); err != nil {
		send(fmt.Sprintf("Couldn't remove the song: %v", err))
		return
	}
	send("Song removed from the queue")
}
