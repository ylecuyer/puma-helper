package cmd

import (
	"os"
	"time"

	log "github.com/Sirupsen/logrus"
	"github.com/spf13/cobra"
)

func init() {
	// Set logger format
	formatter := new(log.TextFormatter)
	formatter.TimestampFormat = time.Stamp
	formatter.FullTimestamp = true
	log.SetFormatter(formatter)
	log.SetOutput(os.Stdout)
}

// RootCmd represents the base command when called without any subcommands
var RootCmd = &cobra.Command{
	Use:   "Puma-Helper",
	Short: "Puma-Helper CLI",
	Long:  "",
	RunE: func(cmd *cobra.Command, args []string) error {
		cmd.Help()

		return nil
	},
}
