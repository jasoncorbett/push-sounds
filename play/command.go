package play

import (
	"github.com/jasoncorbett/push-sounds/libraries"
	"github.com/jasoncorbett/push-sounds/sound"
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

func playSound(c *cli.Context) error {
	lib, err := libraries.NewSoundLibrary(c.String("library-base"))
	if err != nil {
		return err
	}
	soundFile, err := lib.GetRandomFile(c.StringSlice("libraries"))
	if err != nil {
		return err
	}
	soundToPlay, err := sound.NewFromFile(soundFile)
	if err != nil {
		return err
	}
	return soundToPlay.Play()
}
