package bot

import (
	"bytes"
	"context"
	"errors"
	"net/url"
	"os/exec"
	"strings"

	"github.com/pranavtharoor/mc-manager/config"
	"github.com/pranavtharoor/mc-manager/music"
	"google.golang.org/api/option"
	"google.golang.org/api/youtube/v3"
	"gopkg.in/square/go-jose.v2/json"
)

type videoResponse struct {
	ID      string `json:"id"`
	Title   string `json:"fulltitle"`
	Formats []struct {
		Url string `json:"url"`
	} `json:"formats"`
}

func isYoutubeUrl(str string) bool {
	u, err := url.Parse(strings.TrimSpace(str))
	return err == nil && strings.HasSuffix(u.Host, "youtube.com")
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
	var result *youtube.SearchResult
	for _, item := range response.Items {
		if item.Id.Kind == "youtube#video" {
			result = item
			break
		}
	}
	if result == nil {
		return nil, errors.New("couln't find this song")
	}
	return result, nil
}

func getVideo(videoId string) (*videoResponse, error) {
	cmd := exec.Command("youtube-dl", "--skip-download", "--print-json", "--flat-playlist", videoId)
	var out bytes.Buffer
	cmd.Stdout = &out
	err := cmd.Run()
	if err != nil {
		return nil, err
	}
	var video videoResponse
	err = json.NewDecoder(&out).Decode(&video)
	if err != nil {
		return nil, err
	}
	return &video, nil
}

func getSongFromSearchOrUrl(c config.DjConfiguration, searchOrUrl string) (*music.Song, error) {
	var videoId string

	if isYoutubeUrl(searchOrUrl) {
		parts := strings.Split(searchOrUrl, "watch?v=")
		if len(parts) != 2 {
			return nil, errors.New("couldn't play this song")
		}
		videoId = parts[1]
	} else {
		result, err := searchYoutubeForSong(c, searchOrUrl)
		if err != nil {
			return nil, err
		}
		videoId = result.Id.VideoId
	}

	video, err := getVideo(strings.TrimSpace(videoId))
	if err != nil {
		return nil, err
	}

	if len(video.Formats) == 0 {
		return nil, errors.New("couldn't play this song")
	}

	return music.NewSong(video.Title, video.Formats[0].Url), nil
}
