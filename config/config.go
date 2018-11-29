package config

import v "github.com/spf13/viper"

// InitConfig load config from file and/or the environment.
func InitConfig() error {
	v.SetConfigName("puma-helper")
	v.SetConfigType("yaml")
	//v.AddConfigPath("$HOME")
	//v.AddConfigPath("/etc/")
	v.AddConfigPath(".")

	return v.ReadInConfig()
}
