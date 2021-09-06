package play

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
	"github.com/jasoncorbett/push-sounds/libraries"
	"github.com/urfave/cli/v2"
)

var PlayCommand = &cli.Command{
	Name:   "play",
	Usage:  "Play one of the sounds from the library",
	Action: playSound,
	Flags: []cli.Flag{
		&cli.StringSliceFlag{
			Name:    "libraries",
			Aliases: []string{"l"},
			Usage:   "list all the libraries you can pull a sound from",
			Value:   cli.NewStringSlice("default"),
		},
	},
}

func decode(filename string) (s beep.StreamSeekCloser, format beep.Format, err error) {
	audioFile, err := os.Open(filename)
	if err != nil {
		return nil, beep.Format{}, fmt.Errorf("unable to open audiofile: %s", err.Error())
	}
	extension := path.Ext(filename)
	if extension == ".ogg" {
		return vorbis.Decode(audioFile)
	}
	if extension == ".wav" {
		return wav.Decode(audioFile)
	}
	if extension == ".mp3" {
		return mp3.Decode(audioFile)
	}
	if extension == ".flac" {
		return flac.Decode(audioFile)
	}
	return nil, beep.Format{}, fmt.Errorf("unknown audio format '%s'", extension)
}

func playSound(c *cli.Context) error {
	sound, err := libraries.GetSound(c.StringSlice("libraries"))
	if err != nil {
		return err
	}

	fmt.Printf("Playing %s\n", sound)
	stream, format, err := decode(sound)
	if err != nil {
		return fmt.Errorf("unable to decode: %s", err.Error())
	}
	speaker.Init(format.SampleRate, format.SampleRate.N(time.Second/10))
	done := make(chan bool)
	speaker.Play(beep.Seq(stream, beep.Callback(func() {
		done <- true
	})))

	<-done
	time.Sleep(time.Second / 10)
	return nil
}
