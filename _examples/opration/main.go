package main

import (
	"fmt"
	"os"

	"github.com/flyaways/log"
)

func init() {
	//Collector for escape ouput
	log.Collector()
}

func main() {
	operation := log.New(
		log.LogConfig{
			Level:         log.INFO,
			FirstRollover: true,
			Blocking:      false,
			BufferLength:  10240,
			PrefixName:    "operation.log",
			When:          log.Hour,
			BackupCount:   72,
			Format:        log.OprationFormat},
	)

	if operation == nil {
		fmt.Fprintf(os.Stderr, "init logger error\n")
		return
	}

	operation.Debug("Debug")
	operation.Info("Info")
	operation.Warn("Warn")
	operation.Error("Error")
	operation.Fatal("Fatal")
}
