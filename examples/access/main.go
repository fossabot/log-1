package main

import (
	"fmt"
	"os"

	"github.com/flyaways/log"
)

func main() {
	access := log.New(
		log.LogConfig{
			Level:         log.INFO,
			FirstRollover: true,
			Blocking:      false,
			BufferLength:  10240,
			PrefixName:    "access.log",
			When:          log.Hour,
			BackupCount:   72,
			Format:        log.AccessFormat},
	)

	if access == nil {
		fmt.Fprintf(os.Stderr, "init logger error\n")
		return
	}

	access.Info("access info")
}
