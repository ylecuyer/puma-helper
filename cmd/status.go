package cmd

import (
	"errors"

	helper "github.com/dimelo/puma-helper/helper"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

func init() {
	RootCmd.AddCommand(statusCmd)
	setLocalFlags()
}

// statusCmd represents the status command
var statusCmd = &cobra.Command{
	Use:   "status",
	Short: "Command permit to centralize pumactl status metrics in one place",
	Long:  "",
	RunE: func(cmd *cobra.Command, args []string) error {
		if viper.ConfigFileUsed() == "" {
			return nil
		}
		if err := viper.Unmarshal(&helper.CfgFile); err != nil {
			return err
		}
		if err := ensureMandatoryArgs(); err != nil {
			return err
		}
		if err := helper.RunStatus(); err != nil {
			return err
		}

		return nil
	},
}

func setLocalFlags() {
	statusCmd.Flags().StringVarP(&helper.Filter, "filter", "f", "", "Only show applications who match /w given string")
}

func ensureMandatoryArgs() error {
	missing := ""
	err := ""
	for appname, key := range helper.CfgFile.Applications {
		if key.Path == "" {
			missing += "path "
		}
		if len(missing) > 0 {
			err += missing + "arg(s) missing from " + appname + " app\n"
		}
		missing = ""
	}
	if len(err) > 0 {
		return errors.New(err)
	}
	return nil
}
