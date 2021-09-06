package libraries

import (
	"fmt"
)

func GetSound(libraries []string) (string, error) {
	fmt.Printf("Retrieving sound from library list: %s\n", libraries)
	return "sounds/default/push it, push it real good.ogg", nil
}
