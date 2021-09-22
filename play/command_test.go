package play

import (
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/jasoncorbett/push-sounds/libraries"
	"github.com/jasoncorbett/push-sounds/mock_libraries"
	"github.com/jasoncorbett/push-sounds/mock_sound"
	"github.com/jasoncorbett/push-sounds/sound"
	"github.com/urfave/cli/v2"
)

func createApp(libraryLocation string, libraries ...string) *cli.App {
	return &cli.App{
		Flags: []cli.Flag{
			&cli.PathFlag{
				Name:  "library-base",
				Value: libraryLocation,
			},
		},
		Commands: []*cli.Command{
			{
				Name:   "play",
				Action: playSound,
				Flags: []cli.Flag{
					&cli.StringSliceFlag{
						Name:  "libraries",
						Value: cli.NewStringSlice(libraries...),
					},
				},
			},
		},
	}
}

func TestPlayCommand(t *testing.T) {
	randomSoundName := "random-sound"
	expectedBasePath := "base-library-path"
	orig_nsl := libraries.NewSoundLibrary
	orig_nsff := sound.NewFromFile
	defer func() {
		libraries.NewSoundLibrary = orig_nsl
		sound.NewFromFile = orig_nsff
	}()

	// mocks
	m := gomock.NewController(t)
	msl := mock_libraries.NewMockSoundLibrary(m)
	ms := mock_sound.NewMockSound(m)

	var actualSoundName *string
	var actualBasePath *string
	libraries.NewSoundLibrary = func(basePath string) (libraries.SoundLibrary, error) {
		actualBasePath = &basePath
		return msl, nil
	}
	sound.NewFromFile = func(soundFile string) (sound.Sound, error) {
		actualSoundName = &soundFile
		return ms, nil
	}

	msl.
		EXPECT().
		GetRandomFile([]string{"default"}).
		Return(randomSoundName, nil)

	ms.
		EXPECT().
		Play().
		Return(nil)

	app := createApp(expectedBasePath, "default")
	err := app.Run([]string{"test", "play"})
	if err != nil {
		t.Errorf("Recieved error from running fake app, did not expect that: %s", err.Error())
	}
	if *actualBasePath != expectedBasePath {
		t.Fatalf("Expected base path passed to NewSoundLibrary to be '%s', but was '%s'", expectedBasePath, *actualBasePath)
	}
	if *actualSoundName != randomSoundName {
		t.Fatalf("Expected sound name passed to NewSoundFromFile to be '%s', but was '%s'", randomSoundName, *actualSoundName)
	}

}
