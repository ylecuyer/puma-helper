package helper

type Application struct {
	Path           string `mapstructure:"path"`
	Description    string `mapstructure:"description"`
	PumactlPath    string `mapstructure:"pumactlpath"`
	PumastatePath  string `mapstructure:"pumastatepath"`
	ThreadWarn     int    `mapstructure:"thread_warn"`
	ThreadCritical int    `mapstructure:"thread_critical"`
	CPUWarn        int    `mapstructure:"cpu_warn"`
	CPUCritical    int    `mapstructure:"cpu_critical"`
	MemoryWarn     int    `mapstructure:"memory_warn"`
	MemoryCritical int    `mapstructure:"memory_critical"`
}

type Configuration struct {
	Applications map[string]Application `mapstructure:"applications"`
}

var (
	CfgFile Configuration
	Filter  string

	currentApp string
)
