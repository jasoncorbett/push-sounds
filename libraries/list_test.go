package libraries

import (
	"fmt"
	"github.com/golang/mock/gomock"
	"github.com/jasoncorbett/push-sounds/mock_libraries"
	"github.com/urfave/cli/v2"
	"strings"
	"testing"
)

type mockNewSoundLibrary struct {
	Called        bool
	Param         string
	ReturnLibrary SoundLibrary
	ReturnError   error
	original      func(string) (SoundLibrary, error)
}

func (m *mockNewSoundLibrary) Mock() {
	m.original = NewSoundLibrary
	NewSoundLibrary = func(basePath string) (SoundLibrary, error) {
		m.Called = true
		m.Param = basePath
		return m.ReturnLibrary, m.ReturnError
	}
}

func (m *mockNewSoundLibrary) Restore() {
	NewSoundLibrary = m.original
}

func setup(t *testing.T, basePath string) (*gomock.Controller, *mock_libraries.MockSoundLibrary, *cli.App, *mockNewSoundLibrary) {
	m := gomock.NewController(t)
	msl := mock_libraries.NewMockSoundLibrary(m)

	mockNewLibrary := &mockNewSoundLibrary{
		ReturnLibrary: msl,
		ReturnError:   nil,
	}
	app := &cli.App{
		Flags: []cli.Flag{
			&cli.PathFlag{
				Name:  "library-base",
				Value: basePath,
			},
		},
		Commands: []*cli.Command{
			ListCommand,
		},
	}
	return m, msl, app, mockNewLibrary
}

func TestListLibrariesHappyPath(t *testing.T) {
	basePath := "base-path"
	_, msl, app, mockNewLibrary := setup(t, basePath)

	mockNewLibrary.Mock()
	defer mockNewLibrary.Restore()
	msl.
		EXPECT().
		ListLibraries().
		Return([]string{"one", "two"}, nil)

	msl.EXPECT().
		ListFiles("one").
		Return([]string{"a"}, nil)
	msl.EXPECT().
		ListFiles("two").
		Return([]string{"b", "c"}, nil)

	err := app.Run([]string{"push-sounds", "list", "libraries"})

	if err != nil {
		t.Errorf("App should have had no error, but returned %s", err.Error())
	}
	if !mockNewLibrary.Called {
		t.Fatalf("NewSoundLibrary not called from list library command!")
	}
	if mockNewLibrary.Param != basePath {
		t.Fatalf("Expected default base path parameter of NewSoundLibrary to be '%s', but was '%s'", basePath, mockNewLibrary.Param)
	}
}

func TestListLibrariesErrorCreatingSoundLibrary(t *testing.T) {
	basePath := "base-path"
	_, _, app, mockNewLibrary := setup(t, basePath)

	mockNewLibrary.ReturnError = fmt.Errorf("planned testing error")
	mockNewLibrary.ReturnLibrary = nil

	mockNewLibrary.Mock()
	defer mockNewLibrary.Restore()

	err := app.Run([]string{"push-sounds", "list", "libraries"})
	if err == nil {
		t.Fatalf("No error returned from app.Run, but one should have been returned.")
	}
	if !strings.Contains(err.Error(), mockNewLibrary.ReturnError.Error()) {
		t.Fatalf("Expected error from app to contain testing error:\n\tactual error: %s\n\tshould have contained: %s", err.Error(), mockNewLibrary.ReturnError.Error())
	}
}

func TestListLibrariesUnableToListLibraries(t *testing.T) {
	basePath := "base-path"
	_, msl, app, mockNewLibrary := setup(t, basePath)

	mockNewLibrary.Mock()
	defer mockNewLibrary.Restore()

	listLibrariesError := fmt.Errorf("planned testing error")
	msl.
		EXPECT().
		ListLibraries().
		Return([]string{}, listLibrariesError)
	err := app.Run([]string{"push-sounds", "list", "libraries"})
	if err == nil {
		t.Fatalf("No error returned from app.Run, but one should have been returned.")
	}
	if !strings.Contains(err.Error(), listLibrariesError.Error()) {
		t.Fatalf("Expected error from app to contain testing error:\n\tactual error: %s\n\tshould have contained: %s", err.Error(), listLibrariesError.Error())
	}
}

func TestListLibrariesListFilesError(t *testing.T) {
	basePath := "base-path"
	_, msl, app, mockNewLibrary := setup(t, basePath)

	mockNewLibrary.Mock()
	defer mockNewLibrary.Restore()
	msl.
		EXPECT().
		ListLibraries().
		Return([]string{"one"}, nil)

	msl.EXPECT().
		ListFiles("one").
		Return([]string{}, fmt.Errorf("planned testing error"))

	err := app.Run([]string{"push-sounds", "list", "libraries"})
	if err != nil {
		t.Errorf("App should have had no error, but returned %s", err.Error())
	}
}

func TestListLibraryFilesHappyPath(t *testing.T) {
	basePath := "base-path"
	_, msl, app, mockNewLibrary := setup(t, basePath)

	mockNewLibrary.Mock()
	defer mockNewLibrary.Restore()

	msl.EXPECT().
		ListFiles("one").
		Return([]string{"a"}, nil)
	msl.EXPECT().
		ListFiles("two").
		Return([]string{"b", "c"}, nil)

	err := app.Run([]string{"push-sounds", "list", "library", "one", "two"})

	if err != nil {
		t.Errorf("App should have had no error, but returned %s", err.Error())
	}
	if !mockNewLibrary.Called {
		t.Fatalf("NewSoundLibrary not called from list library command!")
	}
	if mockNewLibrary.Param != basePath {
		t.Fatalf("Expected default base path parameter of NewSoundLibrary to be '%s', but was '%s'", basePath, mockNewLibrary.Param)
	}
}

func TestListLibraryErrorCreatingSoundLibrary(t *testing.T) {
	basePath := "base-path"
	_, _, app, mockNewLibrary := setup(t, basePath)

	mockNewLibrary.ReturnError = fmt.Errorf("planned testing error")
	mockNewLibrary.ReturnLibrary = nil

	mockNewLibrary.Mock()
	defer mockNewLibrary.Restore()

	err := app.Run([]string{"push-sounds", "list", "library", "one"})
	if err == nil {
		t.Fatalf("No error returned from app.Run, but one should have been returned.")
	}
	if !strings.Contains(err.Error(), mockNewLibrary.ReturnError.Error()) {
		t.Fatalf("Expected error from app to contain testing error:\n\tactual error: %s\n\tshould have contained: %s", err.Error(), mockNewLibrary.ReturnError.Error())
	}
}

func TestListLibraryErrorListingFiles(t *testing.T) {
	basePath := "base-path"
	_, msl, app, mockNewLibrary := setup(t, basePath)

	mockNewLibrary.Mock()
	defer mockNewLibrary.Restore()

	msl.EXPECT().
		ListFiles("one").
		Return([]string{}, fmt.Errorf("planned testing error"))

	err := app.Run([]string{"push-sounds", "list", "library", "one"})
	if err != nil {
		t.Errorf("App should have had no error, but returned %s", err.Error())
	}
}
