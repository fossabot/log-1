package main

import (
	"fmt"
	"os"

	"github.com/flyaways/log"
)

func main() {
	messageonly := log.New(
		log.LogConfig{
			Level:         log.INFO,
			FirstRollover: true,
			Blocking:      false,
			BufferLength:  10240,
			PrefixName:    "messageonly.log",
			When:          log.Hour,
			BackupCount:   72,
			Format:        log.MessageOnly},
	)

	if messageonly == nil {
		fmt.Fprintf(os.Stderr, "init logger error\n")
		return
	}

	messageonly.Debug("Debug")
	messageonly.Info("Info")
	messageonly.Warn("Warn")
	messageonly.Error("Warn")
}
