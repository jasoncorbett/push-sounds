package sound

type AudioFileType int

type Sound interface {
	Play() error
	Type() AudioFileType
	Location() string
}
