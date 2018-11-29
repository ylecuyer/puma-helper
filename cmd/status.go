package cmd

import (
	helper "github.com/dimelo/puma-helper/helper"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

func init() {
	RootCmd.AddCommand(statusCmd)
}

// statusCmd represents the status command
var statusCmd = &cobra.Command{
	Use:   "status",
	Short: "Give intranet status access in continue",
	Long:  "",
	RunE: func(cmd *cobra.Command, args []string) error {
		if err := viper.Unmarshal(&helper.CfgFile); err != nil {
			return err
		}
		helper.RunStatus()

		return nil
	},
}
