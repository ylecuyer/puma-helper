package main

import (
	log "github.com/Sirupsen/logrus"
	"github.com/dimelo/puma-helper/cmd"
	"github.com/dimelo/puma-helper/config"
)

func main() {
	if err := config.InitConfig(); err != nil {
		log.Error(err)
	}
	if err := cmd.RootCmd.Execute(); err != nil {
		log.Error(err)
	}
}
