package sound

import (
	"math"
	"strings"
	"testing"
)

func listContains(list []AudioFileType, item AudioFileType) bool {
	for _, a := range list {
		if a == item {
			return true
		}
	}
	return false
}

func TestAudioFileType_Name(t *testing.T) {
	for _, audioType := range KnownAudioFileTypes() {
		name := audioType.Name()
		if strings.Contains(strings.ToUpper(name), "UNKNOWN") {
			t.Errorf("known types contains a type that include an unknown name, extension: %s", audioType.Extension())
		}
	}
}

func TestAudioFileType_NameUnknown(t *testing.T) {
	// MaxInt is the Unknown type, which may be handled specially,
	// this test is to handle a type that is not specifically coded in
	audioType := AudioFileType(math.MaxInt - 1)
	if !strings.Contains(strings.ToUpper(audioType.Name()), "UNKNOWN") {
		t.Errorf("The name for an unknown audio type should contain unknown: %s", audioType.Name())
	}
}

func TestAudioFileType_Extension(t *testing.T) {
	for _, audioType := range KnownAudioFileTypes() {
		extension := audioType.Extension()
		if strings.Contains(strings.ToUpper(extension), "UNKNOWN") {
			t.Errorf("known types contains a type that include an unknown extension, name: %s", audioType.Name())
		}
	}
}

func TestAudioFileType_ExtensionUnknown(t *testing.T) {
	// MaxInt is the Unknown type, which may be handled specially,
	// this test is to handle a type that is not specifically coded in
	audioType := AudioFileType(math.MaxInt - 1)
	if !strings.Contains(strings.ToUpper(audioType.Extension()), "UNKNOWN") {
		t.Errorf("The extension for an unknown audio type should contain unknown: %s", audioType.Name())
	}
}

func TestKnownAudioFileTypes(t *testing.T) {
	knownTypes := KnownAudioFileTypes()
	if len(knownTypes) != 4 {
		t.Error("An audio type has been added, but the KnownAudioFileTypesTest has not been adjusted")
	}
	if !listContains(knownTypes, OggVorbisFile) {
		t.Error("KnownAudioFileTypes() missing ogg/vorbis type")
	}
	if !listContains(knownTypes, WavFile) {
		t.Error("KnownAudioFileTypes() missing wav type")
	}
	if !listContains(knownTypes, Mp3File) {
		t.Error("KnownAudioFileTypes() missing mp3 type")
	}
	if !listContains(knownTypes, FlacFile) {
		t.Error("KnownAudioFileTypes() missing flac type")
	}
}
