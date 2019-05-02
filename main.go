package main

import (
	"github.com/dimelo/puma-helper/cmd"
	"github.com/dimelo/puma-helper/config"

	log "github.com/sirupsen/logrus"
)

func main() {
	if err := config.InitConfig(); err != nil {
		log.Warn(err)
	}
	if err := cmd.RootCmd.Execute(); err != nil {
		log.Error(err)
	}
}
