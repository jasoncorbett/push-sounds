package sound

import (
	"fmt"
	"os"
	"path"
	"time"

	"github.com/faiface/beep"
	"github.com/faiface/beep/flac"
	"github.com/faiface/beep/mp3"
	"github.com/faiface/beep/speaker"
	"github.com/faiface/beep/vorbis"
	"github.com/faiface/beep/wav"
)

type beepSound struct {
	Path   string
	stream beep.StreamSeekCloser
	format beep.Format
}

func NewFromFile(soundFile string) (Sound, error) {
	audioFile, err := os.Open(soundFile)
	if err != nil {
		return nil, fmt.Errorf("unable to open audiofile: %s", err.Error())
	}
	extension := path.Ext(soundFile)
	var stream beep.StreamSeekCloser
	var format beep.Format

	switch extension {
	case OggVorbisFile.Extension():
		stream, format, err = vorbis.Decode(audioFile)
	case WavFile.Extension():
		stream, format, err = wav.Decode(audioFile)
	case Mp3File.Extension():
		stream, format, err = mp3.Decode(audioFile)
	case FlacFile.Extension():
		stream, format, err = flac.Decode(audioFile)
	default:
		err = fmt.Errorf("invalid audio file with extension %s", extension)
	}
	if err != nil {
		return nil, fmt.Errorf("unable to decode audio file %s: %s", soundFile, err.Error())
	}
	return &beepSound{
		Path:   soundFile,
		stream: stream,
		format: format,
	}, nil

}

func (bs *beepSound) Play() error {
	err := speaker.Init(bs.format.SampleRate, bs.format.SampleRate.N(time.Second/10))
	if err != nil {
		return fmt.Errorf("unable to initialize audio: %s", err.Error())
	}
	done := make(chan bool)
	speaker.Play(beep.Seq(bs.stream, beep.Callback(func() {
		done <- true
	})))

	<-done
	time.Sleep(time.Second / 10)
	return nil
}

func (bs *beepSound) Location() string {
	return bs.Path
}

func (bs *beepSound) Type() AudioFileType {
	extension := path.Ext(bs.Path)
	for _, aft := range KnownAudioFileTypes() {
		if aft.Extension() == extension {
			return aft
		}
	}
	return UnknownAudioFile
}
