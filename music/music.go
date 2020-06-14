package music

type Song struct {
	Media    string
	Title    string
	Duration *string
	ID       string
}

func NewSong(media, title, id string) *Song {
	song := new(Song)
	song.Media = media
	song.Title = title
	song.ID = id
	return song
}
