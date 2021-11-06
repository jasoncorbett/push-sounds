package libraries

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

var (
	tempLibrary = []string{
		"a/b.wav",
		"c/d.mp3",
		"c/e.ogg",
		"f.txt",
	}
)

func createTempLibrary() (string, error) {
	basePath, err := os.MkdirTemp("", "temp-library-*")
	if err != nil {
		return "", fmt.Errorf("unable to create base library: %s", err.Error())
	}
	for _, path := range tempLibrary {
		parts := strings.Split(path, "/")
		filePart := filepath.Join(basePath, path)
		if len(parts) > 1 {
			dir := filepath.Join(basePath, parts[0])
			if _, err := os.Stat(dir); err != nil {
				os.Mkdir(dir, 0755)
			}
			filePart = filepath.Join(dir, parts[1])
		}
		os.WriteFile(filePart, []byte("test content"), 0666)
	}
	return basePath, nil
}

func removeTempLibrary(basePath string) {
	os.RemoveAll(basePath)
}

func listContains(list []string, item string) bool {
	for _, a := range list {
		if a == item {
			return true
		}
	}
	return false
}

func TestNewSoundLibraryHappyPath(t *testing.T) {
	_, err := NewSoundLibrary(".")
	if err != nil {
		t.Fatalf("error recieved when creating a sound library of '.': %s", err.Error())
	}
}

func TestNewSoundLibraryInvalidPaths(t *testing.T) {
	for _, path := range []string{"does-not-exist", "😊", "://foo"} {
		_, err := NewSoundLibrary(path)
		if err == nil {
			t.Errorf("No error returned for NewSoundLibrary('%s')", path)
		}
	}
}

func TestNewSoundLibraryForFile(t *testing.T) {
	tmpFile, err := os.CreateTemp("", "test-temp")
	if err != nil {
		t.Fatalf("Unable to create tempfile: %s", err.Error())
	}
	defer os.Remove(tmpFile.Name())
	_, err = NewSoundLibrary(tmpFile.Name())
	if err == nil {
		t.Fatalf("NewSoundLibrary('%s') returned no error, but '%s' was a file.", tmpFile.Name(), tmpFile.Name())
	}
}

func TestDirectoryBasedSoundLibrary_ListLibraries(t *testing.T) {
	basePath, err := createTempLibrary()
	defer removeTempLibrary(basePath)
	if err != nil {
		t.Fatalf("Unable to create temporary library: %s", err.Error())
	}
	soundLib, err := NewSoundLibrary(basePath)
	if err != nil {
		t.Fatalf("Error creating sound library for testing: %s", err.Error())
	}
	libraries, err := soundLib.ListLibraries()
	if err != nil {
		t.Fatalf("Error listing libraries from temp library: %s", err.Error())
		return
	}
	if len(libraries) != 2 {
		t.Errorf("ListLibraries return should only contain 'a', and 'c', was: %#v", libraries)
	}
	if !listContains(libraries, "a") || !listContains(libraries, "c") {
		t.Errorf("Libraries should have been 'a' and 'c', was: %#v", libraries)
	}

}

func TestDirectoryBasedSoundLibrary_ListLibrariesErrorDirectoryRemoved(t *testing.T) {
	basePath, err := createTempLibrary()
	// still defer, just in case we error before removing temp library
	defer removeTempLibrary(basePath)
	if err != nil {
		t.Fatalf("Unable to create temporary library: %s", err.Error())
	}
	soundLib, err := NewSoundLibrary(basePath)
	if err != nil {
		t.Fatalf("Error creating sound library for testing: %s", err.Error())
	}
	removeTempLibrary(basePath)
	libraries, err := soundLib.ListLibraries()
	if err == nil {
		t.Error("No error was returned from list libraries call even though directory was removed!")
	}
	if len(libraries) != 0 {
		t.Errorf("No libraries should have been found, instead the following was found: %#v", libraries)
	}
}

func TestDirectoryBasedSoundLibrary_ListFiles(t *testing.T) {
	basePath, err := createTempLibrary()
	defer removeTempLibrary(basePath)
	if err != nil {
		t.Fatalf("Unable to create temporary library: %s", err.Error())
	}
	soundLib, err := NewSoundLibrary(basePath)
	if err != nil {
		t.Fatalf("Error creating sound library for testing: %s", err.Error())
	}
	files, err := soundLib.ListFiles("a")
	if err != nil {
		t.Errorf("Error occurred while listing files in library a: %s", err.Error())
	}
	if len(files) != 1 {
		t.Errorf("Files in library 'a' should have just been 'b.wav', was: %#v", files)
	}
	if !listContains(files, filepath.Join(basePath, "a", "b.wav")) {
		t.Errorf("Files in library 'a' should have just been 'b.wav', was: %#v", files)
	}
	files, err = soundLib.ListFiles("c")
	if err != nil {
		t.Errorf("Error occurred while listing files in library a: %s", err.Error())
	}
	if len(files) != 2 {
		t.Errorf("Files in library 'c' should have been 'd.mp3' and 'e.ogg', was: %#v", files)
	}
	if !listContains(files, filepath.Join(basePath, "c", "d.mp3")) || !listContains(files, filepath.Join(basePath, "c", "e.ogg")) {
		t.Errorf("Files in library 'c' should have been 'd.mp3' and 'e.ogg', was: %#v", files)
	}
}

func TestDirectoryBasedSoundLibrary_ListFilesLibraryDoesNotExist(t *testing.T) {
	basePath, err := createTempLibrary()
	defer removeTempLibrary(basePath)
	if err != nil {
		t.Fatalf("Unable to create temporary library: %s", err.Error())
	}
	soundLib, err := NewSoundLibrary(basePath)
	if err != nil {
		t.Fatalf("Error creating sound library for testing: %s", err.Error())
	}
	files, err := soundLib.ListFiles("doesnotexist")
	if err == nil {
		t.Errorf("Listing a non-existing library did not create and error, files: %#v", files)
	}
	if len(files) != 0 {
		t.Errorf("Listing a non-existing library should not return any files, but did: %#v", files)
	}
}

func TestDirectoryBasedSoundLibrary_ListFilesErrorDirectoryRemoved(t *testing.T) {
	basePath, err := createTempLibrary()
	// still defer, just in case we error before removing temp library
	defer removeTempLibrary(basePath)
	if err != nil {
		t.Fatalf("Unable to create temporary library: %s", err.Error())
	}
	soundLib, err := NewSoundLibrary(basePath)
	if err != nil {
		t.Fatalf("Error creating sound library for testing: %s", err.Error())
	}
	removeTempLibrary(basePath)
	files, err := soundLib.ListFiles("a")
	if err == nil {
		t.Errorf("Listing a removed base-path did not create and error, files: %#v", files)
	}
	if len(files) != 0 {
		t.Errorf("Listing a removed base-path should not return any files, but did: %#v", files)
	}

}

func TestDirectoryBasedSoundLibrary_GetRandomFile(t *testing.T) {
	basePath, err := createTempLibrary()
	defer removeTempLibrary(basePath)
	if err != nil {
		t.Fatalf("Unable to create temporary library: %s", err.Error())
	}
	soundLib, err := NewSoundLibrary(basePath)
	if err != nil {
		t.Fatalf("Error creating sound library for testing: %s", err.Error())
	}
	files, err := soundLib.ListFiles("a")
	if err != nil {
		t.Fatalf("Unable to list files in library during setup for testing: %s", err.Error())
	}
	if len(files) != 1 {
		t.Fatalf("Library 'a' should only contain 1 file, contains: %#v", files)
	}

	randomFile, err := soundLib.GetRandomFile([]string{"a"})
	if err != nil {
		t.Errorf("Error getting random file from library 'a': %s", err.Error())
	}
	if randomFile != files[0] {
		t.Errorf("GetRandomFile(['a']) should have returned %s, instead: %s", files[0], randomFile)
	}
}

func TestDirectoryBasedSoundLibrary_GetRandomFileDoesNotRepeatOften(t *testing.T) {
	basePath, err := createTempLibrary()
	defer removeTempLibrary(basePath)
	if err != nil {
		t.Fatalf("Unable to create temporary library: %s", err.Error())
	}
	soundLib, err := NewSoundLibrary(basePath)
	if err != nil {
		t.Fatalf("Error creating sound library for testing: %s", err.Error())
	}
	allfiles, err := soundLib.ListFiles("a")
	if err != nil {
		t.Fatalf("Unable to list files in library during setup for testing: %s", err.Error())
	}
	files, err := soundLib.ListFiles("c")
	if err != nil {
		t.Fatalf("Unable to list files in library during setup for testing: %s", err.Error())
	}
	allfiles = append(allfiles, files...)
	if len(allfiles) < 3 {
		t.Fatalf("For random file test to work, there should be at least 3 files: %#v", allfiles)
	}
	first, err := soundLib.GetRandomFile([]string{"a", "c"})
	if err != nil {
		t.Fatalf("Error encountered when getting random file: %s", err.Error())
	}
	second, err := soundLib.GetRandomFile([]string{"a", "c"})
	if err != nil {
		t.Fatalf("Error encountered when getting random file: %s", err.Error())
	}
	third, err := soundLib.GetRandomFile([]string{"a", "c"})
	if err != nil {
		t.Fatalf("Error encountered when getting random file: %s", err.Error())
	}
	if first == second && second == third {
		t.Errorf("Getting 3 random files in a row should not have produced the same result: %s", first)
	}
	if !listContains(allfiles, first) || !listContains(allfiles, second) || !listContains(allfiles, third) {
		t.Errorf("Both the first, second, and third response should have been in the list of all files.  First: %s; Second: %s; Third: %s; All Files: %#v", first, second, third, allfiles)
	}
}

func TestDirectoryBasedSoundLibrary_GetRandomFileNonExistingLibrary(t *testing.T) {
	basePath, err := createTempLibrary()
	defer removeTempLibrary(basePath)
	if err != nil {
		t.Fatalf("Unable to create temporary library: %s", err.Error())
	}
	soundLib, err := NewSoundLibrary(basePath)
	if err != nil {
		t.Fatalf("Error creating sound library for testing: %s", err.Error())
	}
	randomFile, err := soundLib.GetRandomFile([]string{"doesnotexist"})
	if err == nil {
		t.Errorf("GetRandomFile(['doesnotexist']) should have returned and error, but didn't.  RandomFile: %s", randomFile)
	}
	if randomFile != "" {
		t.Errorf("GetRandomFile(['doesnotexist']) should have returned and empty string, instead: %s", randomFile)
	}
}

func TestDirectoryBasedSoundLibrary_GetRandomFileIncludingNonExistingLibrary(t *testing.T) {
	basePath, err := createTempLibrary()
	defer removeTempLibrary(basePath)
	if err != nil {
		t.Fatalf("Unable to create temporary library: %s", err.Error())
	}
	soundLib, err := NewSoundLibrary(basePath)
	if err != nil {
		t.Fatalf("Error creating sound library for testing: %s", err.Error())
	}
	files, err := soundLib.ListFiles("a")
	if err != nil {
		t.Fatalf("Unable to list files in library during setup for testing: %s", err.Error())
	}
	if len(files) != 1 {
		t.Fatalf("Library 'a' should only contain 1 file, contains: %#v", files)
	}
	randomFile, err := soundLib.GetRandomFile([]string{"a", "doesnotexist"})
	if err != nil {
		t.Fatalf("Getting a random file from a valid library and a non-existing one should not return error but did: %s", err.Error())
	}
	if randomFile != files[0] {
		t.Errorf("GetRandomFile(['a', 'doesnotexist']) should have returned %s, instead: %s", files[0], randomFile)
	}
}

func TestDirectoryBasedSoundLibrary_GetRandomFileRemovedBasePath(t *testing.T) {
	basePath, err := createTempLibrary()
	// still defer, just in case we error before removing temp library
	defer removeTempLibrary(basePath)
	if err != nil {
		t.Fatalf("Unable to create temporary library: %s", err.Error())
	}
	soundLib, err := NewSoundLibrary(basePath)
	if err != nil {
		t.Fatalf("Error creating sound library for testing: %s", err.Error())
	}
	removeTempLibrary(basePath)
	randomFile, err := soundLib.GetRandomFile([]string{"a", "c"})
	if err == nil {
		t.Errorf("Getting a random file from a removed base-path did not create and error, randomFile: %s", randomFile)
	}
	if randomFile != "" {
		t.Errorf("Getting a random file from a removed base-path should have returned an empty string, but returned: %s", randomFile)
	}
}
