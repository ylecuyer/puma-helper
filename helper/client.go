package helper

type Application struct {
	Path          string `mapstructure:"path"`
	Description   string `mapstructure:"description"`
	PumactlPath   string `mapstructure:"pumactlpath"`
	PumastatePath string `mapstructure:"pumastatepath"`
}

type Configuration struct {
	Applications map[string]Application `mapstructure:"applications"`
}

var (
	CfgFile Configuration
)
