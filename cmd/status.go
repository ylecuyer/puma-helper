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
	Short: "Command permit to centralize pumactl status metrics in one place",
	Long:  "",
	RunE: func(cmd *cobra.Command, args []string) error {
		if err := viper.Unmarshal(&helper.CfgFile); err != nil {
			return err
		}
		if err := helper.RunStatus(); err != nil {
			return err
		}

		return nil
	},
}
