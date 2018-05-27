package main

import (
	"fmt"
	"os"
	"os/signal"

	"github.com/flyaways/log"
)

func init() {
	//Collector for escape ouput
	log.Collector()
}

func main() {
	op := log.New(
		log.LogConfig{
			Level:         log.INFO,
			FirstRollover: true,
			Blocking:      false,
			BufferLength:  10240,
			PrefixName:    "op.log",
			When:          log.Hour,
			BackupCount:   72,
			Format:        log.OpFormat}, //log.AccessFormat
	)

	if op == nil {
		fmt.Fprintf(os.Stderr, "init logger error\n")
		return
	}

	op.Debug("Go is an open source project developed by a team at Google and many contributors from the open source community.")
	op.Info("Go is distributed under a BSD-style license.")
	op.Warn("A low traffic mailing list for important announcements, such as new releases.")
	op.Error("We encourage all Go users to subscribe to golang-announce.")
	//op.Fatal("Go is an open source programming language that makes it easy to build simple, reliable, and efficient software.")

	signals := make(chan os.Signal, 1)
	signal.Notify(signals, os.Interrupt)
	<-signals
}
