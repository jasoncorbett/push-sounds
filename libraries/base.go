package libraries

import (
	"fmt"
	"io/ioutil"
	"math/rand"
	"os"
	"path/filepath"
	"time"
)

type SoundLibrary interface {
	ListLibraries() ([]string, error)
	GetRandomFile(from []string) (string, error)
	ListFiles(library string) ([]string, error)
}

type directoryBasedSoundLibrary struct {
	basePath string
}

func NewSoundLibrary(basePath string) (SoundLibrary, error) {
	stat, err := os.Stat(basePath)
	if err != nil && os.IsNotExist(err) {
		return &directoryBasedSoundLibrary{}, fmt.Errorf("library base %s does not exist", basePath)
	}
	if err != nil {
		return &directoryBasedSoundLibrary{}, fmt.Errorf("unable to initialize library: %s", err.Error())
	}
	if !stat.IsDir() {
		return &directoryBasedSoundLibrary{}, fmt.Errorf("library base %s is not a directory", basePath)
	}
	return &directoryBasedSoundLibrary{
		basePath: basePath,
	}, nil
}

func (l *directoryBasedSoundLibrary) ListLibraries() ([]string, error) {
	files, err := ioutil.ReadDir(l.basePath)
	if err != nil {
		return []string{}, fmt.Errorf("unable to read %s: %s", l.basePath, err.Error())
	}
	libraries := []string{}
	for _, library := range files {
		if library.IsDir() {
			libraries = append(libraries, library.Name())
		}
	}
	return libraries, nil
}

func (l *directoryBasedSoundLibrary) GetRandomFile(from []string) (string, error) {
	files := []string{}
	for _, library := range from {
		libraryFiles, err := l.ListFiles(library)
		if err == nil {
			files = append(files, libraryFiles...)
		}
	}
	if len(files) == 0 {
		return "", fmt.Errorf("no files available in %v", from)
	}
	rand.Seed(time.Now().Unix())
	randIndex := rand.Intn(len(files))
	return files[randIndex], nil
}

func (l *directoryBasedSoundLibrary) ListFiles(library string) ([]string, error) {
	libraryPath := filepath.Join(l.basePath, library)
	files, err := ioutil.ReadDir(libraryPath)
	if err != nil {
		return []string{}, fmt.Errorf("unable to read %s: %s", libraryPath, err)
	}
	libraryFiles := []string{}
	for _, libraryFile := range files {
		if !libraryFile.IsDir() {
			libraryFiles = append(libraryFiles, filepath.Join(libraryPath, libraryFile.Name()))
		}
	}
	return libraryFiles, nil

}
