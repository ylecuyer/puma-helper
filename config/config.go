package config

import (
	"os"
	"path/filepath"

	v "github.com/spf13/viper"
)

// InitConfig load config from file and/or the environment.
func InitConfig() error {
	v.SetConfigName(".puma-helper")
	v.SetConfigType("yaml")
	v.AddConfigPath("$HOME")
	v.AddConfigPath(".")

	dir, err := filepath.Abs(filepath.Dir(os.Args[0]))
	if err != nil {
		return err
	}
	v.AddConfigPath(dir)

	return v.ReadInConfig()
}
