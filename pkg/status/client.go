package status

// PumaHelperCfg point to struct who contain all options from puma-helper.yml file
type PumaHelperCfg struct {
	Applications map[string]PumaHelperCfgData `mapstructure:"applications"`
}

// PumaHelperCfgData contain all options per app from puma-helper.yml file
type PumaHelperCfgData struct {
	Path           string   `mapstructure:"path"`
	Description    string   `mapstructure:"description"`
	PumastatePaths []string `mapstructure:"pumastatepaths"`
	ThreadWarn     int      `mapstructure:"thread_warn"`
	ThreadCritical int      `mapstructure:"thread_critical"`
	CPUWarn        int      `mapstructure:"cpu_warn"`
	CPUCritical    int      `mapstructure:"cpu_critical"`
	MemoryWarn     int      `mapstructure:"memory_warn"`
	MemoryCritical int      `mapstructure:"memory_critical"`
}

const (
	appGroupName  string = "/current"
	pumastatePath string = appGroupName + "/tmp/pids/puma.state"
)

var (
	// CfgFile point to struct who contain all options from puma-helper.yml file
	CfgFile PumaHelperCfg
	// Filter CLI status command option
	Filter string
	// JSONOutput CLI status command option
	JSONOutput bool
	// ExpandDetails CLI status command option
	ExpandDetails bool

	currentApp string

	retrieveStatusError []string
)
