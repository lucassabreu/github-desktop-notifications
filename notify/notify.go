package notify

import (
	"fmt"

	"github.com/gen2brain/beeep"
)

// Notify the contents
func Notify(reason, message, url string) error {
	return beeep.Notify(
		fmt.Sprint("GitHub ", reason),
		fmt.Sprint(message, "\n", url),
		"assets/information.png",
	)
}
