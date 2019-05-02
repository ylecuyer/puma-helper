package cmd

import (
	"os"
	"time"

	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

func init() {
	formatter := new(log.TextFormatter)
	formatter.TimestampFormat = time.Stamp
	formatter.FullTimestamp = true
	log.SetFormatter(formatter)
	log.SetOutput(os.Stdout)
}

// RootCmd represents the base command when called without any subcommands
var RootCmd = &cobra.Command{
	Use:   "Puma-Helper",
	Short: "Puma-Helper CLI aims to implement missing centralized and human readeable features from puma unix socket in one place.",
	Long:  "",
	RunE: func(cmd *cobra.Command, args []string) error {
		cmd.Help()

		return nil
	},
}
