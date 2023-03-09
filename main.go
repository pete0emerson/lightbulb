package main

import (
	"github.com/pete0emerson/lightbulb/cmd"
	log "github.com/sirupsen/logrus"
)

func init() {
	log.SetFormatter(&log.TextFormatter{
		FullTimestamp: true,
	})
	log.SetLevel(log.DebugLevel)
}

func main() {
	cmd.Execute()
}
