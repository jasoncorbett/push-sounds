package libraries

import (
	"os"
	"path/filepath"
)

var (
	GetUserConfigBase = os.UserConfigDir
)

func GetLocationDefault() string {
	configDir, err := GetUserConfigBase()
	if err != nil {
		configDir = "."
	}
	return filepath.Join(configDir, "push-sounds")
}
