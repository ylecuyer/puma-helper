package helper

type Application struct {
	Path          string `mapstructure:"path"`
	State         string `mapstructure:"state"`
	PumactlPath   string `mapstructure:"pumactlpath"`
	PumastatePath string `mapstructure:"pumastatepath"`
}

type Configuration struct {
	Applications map[string]Application `mapstructure:"applications"`
}

var (
	CfgFile Configuration
)
