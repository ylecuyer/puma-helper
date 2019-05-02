package cmd

import (
	"errors"

	status "github.com/dimelo/puma-helper/pkg/status"
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
	Short: "Centralize puma unix socket status metrics in one place",
	Long:  "",
	RunE: func(cmd *cobra.Command, args []string) error {
		if viper.ConfigFileUsed() == "" {
			return nil
		}
		if err := viper.Unmarshal(&status.CfgFile); err != nil {
			return err
		}
		if err := ensureArgsValidity(); err != nil {
			return err
		}
		if err := status.RunStatus(); err != nil {
			return err
		}

		return nil
	},
}

func setLocalFlags() {
	statusCmd.Flags().StringVarP(&status.Filter, "filter", "f", "", "Only show applications who match /w given string")
	statusCmd.Flags().BoolVarP(&status.JSONOutput, "json", "j", false, "Return JSON object who contains all informations")
	statusCmd.Flags().BoolVarP(&status.ExpandDetails, "details", "d", false, "Show more details about apps and workers")
}

func ensureArgsValidity() error {
	err := ""
	finalerr := ""
	for appname, key := range status.CfgFile.Applications {

		if key.Path == "" {
			err += "path missing, "
		}

		if (key.ThreadWarn < 0 || key.ThreadWarn > 100) && key.ThreadWarn != 0 {
			err += "thread warn % is incorrect, "
		}

		if (key.ThreadCritical < 0 || key.ThreadCritical > 100) && key.ThreadCritical != 0 {
			err += "thread critical % is incorrect, "
		}

		if (key.CPUWarn < 0 || key.CPUWarn > 100) && key.CPUWarn != 0 {
			err += "CPU warn % is incorrect, "
		}

		if (key.CPUCritical < 0 || key.CPUCritical > 100) && key.CPUCritical != 0 {
			err += "CPU critical % is incorrect, "
		}

		if key.MemoryWarn < 0 && key.MemoryWarn != 0 {
			err += "memory warn % is incorrect, "
		}

		if key.MemoryWarn < 0 && key.MemoryWarn != 0 {
			err += "memory critical % is incorrect, "
		}

		if len(err) > 0 {
			finalerr += err[:len(err)-2] + " for " + appname + " configuration\n"
		}

		err = ""
	}

	if len(finalerr) > 0 {
		return errors.New(finalerr)
	}

	return nil
}
