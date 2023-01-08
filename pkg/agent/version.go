package agent

import "fmt"

var Version string

func GetVersion() string {
	return fmt.Sprintf("Version: %s", Version)
}
