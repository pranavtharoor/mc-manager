package music

import (
	"errors"
	"time"

	"github.com/bwmarrin/discordgo"
	"github.com/jonas747/dca"
)

type Stereo struct {
	getVoiceConnection func() *discordgo.VoiceConnection
	streamingSession   *dca.StreamingSession
	encodeSession      *dca.EncodeSession
	done               chan error
	queue              []*Song
	nowPlaying         *Song
}

func NewStereo(getVoiceConnection func() *discordgo.VoiceConnection) *Stereo {
	stereo := &Stereo{getVoiceConnection: getVoiceConnection, queue: []*Song{}}
	go stereo.startAutoPlayer()
	return stereo
}

func (s *Stereo) IsConnected() bool {
	return s.getVoiceConnection() != nil
}

func (s *Stereo) Disconnect() {
	if s.IsConnected() {
		s.getVoiceConnection().Disconnect()
	}
}

func (s *Stereo) IsStreaming() bool {
	if s.streamingSession == nil {
		return false
	}
	isFinished, _ := s.streamingSession.Finished()
	return !isFinished
}

func (s *Stereo) StartStreaming(song *Song) error {
	if !s.IsConnected() {
		return errors.New("tried playing without a voice connection")
	}

	s.AddToQueue(song, 0)

	s.Next()

	return nil
}

func (s *Stereo) AddToQueue(song *Song, index int) {
	i := index
	if i < 0 {
		i = 0
	}
	if len(s.queue) == 0 || i >= len(s.queue) {
		s.queue = append(s.queue, song)
		return
	}
	s.queue = append(s.queue[:i+1], s.queue[i:]...)
	s.queue[i] = song
}

func (s *Stereo) RemoveFromQueue(index int) error {
	i := index
	if i < 0 || len(s.queue)-1 < index {
		return errors.New("index out of bounds")
	}
	s.queue = append(s.queue[:index], s.queue[index+1:]...)
	return nil
}

func (s *Stereo) playNext() error {
	if len(s.queue) == 0 {
		return nil
	}
	song := s.queue[0]
	s.queue = s.queue[1:]
	options := dca.StdEncodeOptions
	options.RawOutput = true
	options.Application = "lowdelay"

	encodeSession, err := dca.EncodeFile(song.Url, options)
	if err != nil {
		return err
	}

	s.nowPlaying = song
	s.encodeSession = encodeSession
	s.done = make(chan error)
	s.streamingSession = dca.NewStream(encodeSession, s.getVoiceConnection(), s.done)
	return nil
}

func (s *Stereo) Next() {
	if s.IsStreaming() {
		s.done <- nil
	}
}

func (s *Stereo) cleanup() {
	if s.encodeSession != nil {
		s.encodeSession.Cleanup()
	}
	s.done = nil
	s.streamingSession = nil
	s.encodeSession = nil
	s.nowPlaying = nil
}

func (s *Stereo) ClearQueue() {
	s.queue = []*Song{}
}

func (s *Stereo) GetQueue() []*Song {
	return s.queue
}

func (s *Stereo) GetNowPlaying() *Song {
	return s.nowPlaying
}

func (s *Stereo) startAutoPlayer() {
	for {
		time.Sleep(time.Second) // TODO: fix concurreny issue if possible
		if !s.IsConnected() {
			continue
		}
		time.Sleep(time.Second) // TODO: fix concurreny issue if possible
		s.playNext()
		if s.done == nil {
			continue
		}
		err := <-s.done
		if err != nil {
			println(err.Error())
		}
		s.cleanup()
	}
}
