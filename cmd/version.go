package cmd

import (
	"fmt"

	. "github.com/logrusorgru/aurora"
	"github.com/spf13/cobra"

	version "github.com/dimelo/puma-helper/pkg/version"
)

const (
	currentVersion = "v" + version.Version
)

func init() {
	RootCmd.AddCommand(versionCmd)
}

// versionCmd represents the init command
var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Print actual version",
	Long:  "",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Printf("Current version: %s\n", Gray(currentVersion))
	},
}
