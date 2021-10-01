package bot

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"os/exec"

	"github.com/bwmarrin/discordgo"
	"github.com/pranavtharoor/mc-manager/config"
	"github.com/pranavtharoor/mc-manager/music"
	"google.golang.org/api/option"
	"google.golang.org/api/youtube/v3"
	"gopkg.in/square/go-jose.v2/json"
)

type videoResponse struct {
	Formats []struct {
		Url string `json:"url"`
	} `json:"formats"`
}

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

func djPlay(c config.DjConfiguration, stereo *music.Stereo, song string, send func(string)) {
	if !stereo.IsConnected() {
		send("I'm not in any voice channels")
		return
	}

	video, err := searchYoutubeForSong(c, song)
	if err != nil {
		send(err.Error())
		return
	}

	url, err := getDownloadUrl(video.Id.VideoId)
	if err != nil {
		send(err.Error())
		return
	}

	send("Playing " + video.Snippet.Title)
	err = stereo.StartStreaming(music.NewSong(video.Snippet.Title, url))
	if err != nil {
		send(err.Error())
	}
}

func searchYoutubeForSong(c config.DjConfiguration, search string) (*youtube.SearchResult, error) {
	developerKey := c.YoutubeToken
	ctx := context.Background()
	service, err := youtube.NewService(ctx, option.WithAPIKey(developerKey))
	if err != nil {
		return nil, err
	}
	call := service.Search.List([]string{"id", "snippet"}).Q(search + " lyrics").MaxResults(10)
	response, err := call.Do()
	if err != nil {
		return nil, err
	}
	var video *youtube.SearchResult
	for _, item := range response.Items {
		if item.Id.Kind == "youtube#video" {
			video = item
			break
		}
	}
	if video == nil {
		return nil, errors.New("couln't find this song")
	}
	return video, nil
}

func getDownloadUrl(videoId string) (string, error) {
	cmd := exec.Command("youtube-dl", "--skip-download", "--print-json", "--flat-playlist", videoId)
	var out bytes.Buffer
	cmd.Stdout = &out
	err := cmd.Run()
	if err != nil {
		return "", err
	}
	var videoInfo videoResponse
	err = json.NewDecoder(&out).Decode(&videoInfo)
	if err != nil {
		return "", err
	}
	if len(videoInfo.Formats) == 0 {
		return "", errors.New("couldn't play this song")
	}
	return videoInfo.Formats[0].Url, nil
}

func djList(stereo *music.Stereo) string {
	songs := stereo.GetQueue()
	output := "Queue:\n"
	for i, song := range songs {
		output += fmt.Sprintf("\n%d. %s", i+1, song.Title)
	}
	return output
}

func djAdd(c config.DjConfiguration, stereo *music.Stereo, song string, send func(string), index int) {
	video, err := searchYoutubeForSong(c, song)
	if err != nil {
		send(err.Error())
		return
	}

	url, err := getDownloadUrl(video.Id.VideoId)
	if err != nil {
		send(err.Error())
		return
	}

	send("Adding " + video.Snippet.Title)
	stereo.AddToQueue(music.NewSong(video.Snippet.Title, url), index)
}

func djRemove(stereo *music.Stereo, send func(string), index int) {
	if err := stereo.RemoveFromQueue(index); err != nil {
		send(fmt.Sprintf("Couldn't remove the song: %v", err))
		return
	}
	send("Song removed from the queue")
}
