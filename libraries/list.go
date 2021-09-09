package libraries

import (
	"fmt"
	"path/filepath"
	"strings"

	"github.com/urfave/cli/v2"
)

var ListCommand = &cli.Command{
	Name:  "list",
	Usage: "list libraries or audio files in those libraries",
	Subcommands: []*cli.Command{
		{
			Name:   "libraries",
			Usage:  "List all libraries",
			Action: listLibraries,
		},
		{
			Name:      "library",
			Usage:     "List the files in a library (use additional arguments to for library to list)",
			Action:    listFilesInLibrary,
			ArgsUsage: "*<library name>",
		},
	},
}

func listLibraries(c *cli.Context) error {
	lib, err := NewSoundLibrary(c.String("library-base"))
	if err != nil {
		return fmt.Errorf("unable to initialize sound library: %s", err.Error())
	}
	libraries, err := lib.ListLibraries()
	if err != nil {
		return fmt.Errorf("problem listing existing libraries: %s", err.Error())
	}
	if len(libraries) == 0 {
		fmt.Println("No libraries found")
		return nil
	}
	fmt.Printf("%-15s %s\n", "Name", "Files")
	fmt.Printf("%s %s\n", strings.Repeat("-", 15), strings.Repeat("-", len("Files")))
	for _, library := range libraries {
		numOfLibFiles := 0
		libFiles, err := lib.ListFiles(library)
		if err == nil {
			numOfLibFiles = len(libFiles)
		}
		fmt.Printf("%-15s %5d\n", library, numOfLibFiles)
	}
	return nil
}

func listFilesInLibrary(c *cli.Context) error {
	if c.Args().Len() == 0 {
		return fmt.Errorf("required library name to list")
	}
	lib, err := NewSoundLibrary(c.String("library-base"))
	if err != nil {
		return fmt.Errorf("unable to initialize sound library: %s", err.Error())
	}
	for _, libraryName := range c.Args().Slice() {
		fmt.Printf("%s\n%s\n", libraryName, strings.Repeat("=", len(libraryName)))
		files, err := lib.ListFiles(libraryName)
		if err == nil {
			if len(files) == 0 {
				fmt.Printf("No files found\n\n")
			} else {
				for _, file := range files {
					fmt.Println(filepath.Base(file))
				}
				fmt.Println()
			}
		} else {
			fmt.Printf("Error listing library files: %s\n\n", err.Error())
		}
	}
	return nil
}
