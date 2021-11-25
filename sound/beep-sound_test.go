package sound

import (
	"fmt"
	"testing"
)

func TestBeepSound_PlayInvalid(t *testing.T) {
	sound := beepSound{}
	err := sound.Play()
	if err == nil {
		t.Error("Tried creating an invalid sound, then playing, but it didn't create an error")
	}
}

func TestBeepSound_Location(t *testing.T) {
	sound := beepSound{Path: "some path"}
	location := sound.Location()
	if location != "some path" {
		t.Errorf("sound location should have been 'some path' but was: %#v", location)
	}
}

func TestBeepSound_Type(t *testing.T) {
	for _, aft := range KnownAudioFileTypes() {
		sound := beepSound{Path: fmt.Sprintf("sound%s", aft.Extension())}
		if aft != sound.Type() {
			t.Errorf("sound file '%s' identified as type %s instead of %s", sound.Path, sound.Type().Name(), aft.Name())
		}
	}
}

func TestNewFromFile(t *testing.T) {

}

const (
	MP3_BASE64 = ""
)
