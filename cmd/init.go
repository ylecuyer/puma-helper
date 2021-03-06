package cmd

import (
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	input "github.com/tcnksm/go-input"
	yaml "gopkg.in/yaml.v2"

	config "github.com/dimelo/puma-helper/config"
	status "github.com/dimelo/puma-helper/pkg/status"
)

const (
	globingExpression string = "/home/*/current/tmp/pids/puma*state"
)

func init() {
	RootCmd.AddCommand(initCmd)
}

// initCmd represents the init command
var initCmd = &cobra.Command{
	Use:   "init",
	Short: "Initialize configuration file if it doesn't exist",
	Long:  "",
	RunE: func(cmd *cobra.Command, args []string) error {

		ui := &input.UI{
			Writer: os.Stdout,
			Reader: os.Stdin,
		}

		if viper.ConfigFileUsed() != "" {
			ret, err := ui.Ask("A configuration file already exist, do you want erase? (y/n, required)",
				&input.Options{
					Required:  true,
					HideOrder: true,
					ValidateFunc: func(s string) error {
						if s != "y" && s != "n" {
							return errors.New("Must be y or n")
						}
						return nil
					},
				})
			if err != nil || ret == "n" {
				return err
			}
		}

		gbool, err := ui.Ask("Would you try to use globing to init configuration? (y/n, required)", &input.Options{
			Required:  true,
			HideOrder: true,
			ValidateFunc: func(s string) error {
				if s != "y" && s != "n" {
					return errors.New("Must be y or n")
				}
				return nil
			},
		})
		if err != nil {
			return err
		}

		if gbool == "y" {
			files, err := filepath.Glob("/home/*/current/tmp/pids/puma*state")
			if err != nil {
				return nil
			}
			return buildStructGlobing(files)
		}

		appname, err := ui.Ask("What's your app name? (string, required)", &input.Options{
			Required:  true,
			HideOrder: true,
			ValidateFunc: func(s string) error {
				if s == "" {
					return errors.New("Must be not empty")
				}
				return nil
			},
		})
		if err != nil {
			return err
		}

		apppath, err := ui.Ask("What's absolute path to your puma app? (string, required)", &input.Options{
			Required:  true,
			HideOrder: true,
			ValidateFunc: func(s string) error {
				if s == "" {
					return errors.New("Must be not empty")
				}
				return nil
			},
		})
		if err != nil {
			return err
		}

		pumastatepath, err := ui.Ask("What's absolute path to your puma state file? (string, optionnal)", &input.Options{
			Required:  false,
			HideOrder: true,
		})
		if err != nil {
			return err
		}

		if err := buildAndWriteConfigFile(appname, apppath, pumastatepath); err != nil {
			return err
		}

		return nil
	},
}

func buildStructGlobing(files []string) error {
	cfgdata := make(map[string]status.PumaHelperCfgData)

	for fid := range files {
		fi, err := os.Lstat(files[fid])
		if err != nil {
			return err
		}

		if fi.Mode()&os.ModeSymlink != 0 {
			continue
		}

		cutpath := strings.Split(files[fid], "/")
		if key, ok := cfgdata[cutpath[2]]; ok {
			key.PumastatePaths = append(key.PumastatePaths, files[fid])
			cfgdata[cutpath[2]] = key
		} else {
			cfgdata[cutpath[2]] = status.PumaHelperCfgData{
				Path:           "/home/" + cutpath[2],
				PumastatePaths: []string{files[fid]},
			}
		}
	}

	return marshalAndWriteConfigFile(
		status.PumaHelperCfg{
			Applications: cfgdata,
		})
}

func buildAndWriteConfigFile(appname, apppath, pumastatepath string) error {
	cfgdata := make(map[string]status.PumaHelperCfgData)

	cfgdata[appname] = status.PumaHelperCfgData{
		Path:           apppath,
		PumastatePaths: []string{pumastatepath},
	}

	return marshalAndWriteConfigFile(
		status.PumaHelperCfg{
			Applications: cfgdata,
		})
}

func marshalAndWriteConfigFile(cfg status.PumaHelperCfg) error {
	d, err := yaml.Marshal(&cfg)
	if err != nil {
		return err
	}

	if os.Getenv("HOME") == "" {
		return errors.New("$HOME need to be set")
	}

	fname := config.CfgFileName + "." + config.CfgFileExt
	fnamepath := os.Getenv("HOME") + "/" + fname

	if err := ioutil.WriteFile(fnamepath, d, 0644); err != nil {
		return err
	}

	fmt.Printf("Congratulations! %s is now created -> %s\n", fname, fnamepath)
	fmt.Println("You can now edit it directly if you wan to add more apps or specific options.")

	return nil
}
