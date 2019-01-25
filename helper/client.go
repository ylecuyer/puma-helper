package helper

type pumaHelperCfg struct {
	Applications map[string]struct {
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
	} `mapstructure:"applications"`
}

var (
	// CfgFile point to struct who contain all options from puma-helper.yml file
	CfgFile pumaHelperCfg
	// Filter CLI status command option
	Filter string
	// JSONOutput CLI status command option
	JSONOutput bool

	currentApp string
)
