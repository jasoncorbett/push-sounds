package main

import (
	"log"
	"os"

	"github.com/jasoncorbett/push-sounds/play"
	"github.com/urfave/cli/v2"
)

func main() {
	app := &cli.App{
		Usage: "Play sounds when you git push",
		Commands: []*cli.Command{
			play.PlayCommand,
		},
	}

	err := app.Run(os.Args)
	if err != nil {
		log.Fatal(err)
	}
}
