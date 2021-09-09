package main

import (
	"log"
	"os"

	"github.com/jasoncorbett/push-sounds/libraries"
	"github.com/jasoncorbett/push-sounds/play"
	"github.com/urfave/cli/v2"
)

func main() {
	app := &cli.App{
		Usage: "Play sounds when you git push",
		Flags: []cli.Flag{
			&cli.PathFlag{
				Name:    "library-base",
				Aliases: []string{"lib", "l"},
				Value:   libraries.GetLocationDefault(),
				Usage:   "The base directory for libraries.  Each library will be a directory under this folder.",
				EnvVars: []string{"PUSH_SOUNDS_LIBRARY"},
			},
		},
		Commands: []*cli.Command{
			play.PlayCommand,
			libraries.ListCommand,
		},
	}

	err := app.Run(os.Args)
	if err != nil {
		log.Fatal(err)
	}
}
