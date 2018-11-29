package helper

type Application struct {
	Path  string `mapstructure:"path"`
	State string `mapstructure:"state"`
}

type Configuration struct {
	Applications map[string]Application `mapstructure:"applications"`
}

var (
	CfgFile Configuration
)
