package libraries

import (
	"os"
	"path/filepath"
)

func GetLocationDefault() string {
	configDir, err := os.UserConfigDir()
	if err != nil {
		configDir = "."
	}
	return filepath.Join(configDir, "push-sounds")
}
