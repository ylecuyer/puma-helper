package config

import (
	"os"
	"path/filepath"

	v "github.com/spf13/viper"
)

const (
	//CfgFileName represents the config file name
	CfgFileName string = ".puma-helper"
	//CfgFileExt represents the config file extension
	CfgFileExt string = "yaml"
)

// InitConfig load config from file and/or the environment.
func InitConfig() error {
	v.SetConfigName(CfgFileName)
	v.SetConfigType(CfgFileExt)
	v.AddConfigPath("$HOME")
	v.AddConfigPath(".")

	dir, err := filepath.Abs(filepath.Dir(os.Args[0]))
	if err != nil {
		return err
	}
	v.AddConfigPath(dir)

	return v.ReadInConfig()
}
