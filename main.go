package main

import (
	"log"
	"os"

	"github.com/urfave/cli/v2"
)

func main() {
	app := &cli.App{
		Usage: "play sounds when you git push",
	}

	err := app.Run(os.Args)
	if err != nil {
		log.Fatal(err)
	}
}
