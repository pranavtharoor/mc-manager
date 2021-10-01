package music

type Song struct {
	Title string
	Url   string
}

func NewSong(title string, url string) *Song {
	return &Song{Title: title, Url: url}
}
