package status

type pumaStatus struct {
	BootedWorkers int `json:"booted_workers"`
	OldWorkers    int `json:"old_workers"`
	Phase         int `json:"phase"`
	MainPid       int `json:"main_pid"`
	WorkerStatus  []struct {
		Booted      bool   `json:"booted"`
		Index       int    `json:"index"`
		LastCheckin string `json:"last_checkin"`
		LastStatus  struct {
			Backlog      int `json:"backlog"`
			MaxThreads   int `json:"max_threads"`
			PoolCapacity int `json:"pool_capacity"`
			Running      int `json:"running"`
		} `json:"last_status"`
		Phase int `json:"phase"`
		Pid   int `json:"pid"`
	} `json:"worker_status"`
	Workers int `json:"workers"`
}

type pumaStatusFinalOutput struct {
	ApplicationsCount int                     `json:"applications_count"`
	Application       []pumaStatusApplication `json:"application"`
}

type pumaStatusApplication struct {
	Name           string                 `json:"name"`
	Description    string                 `json:"description"`
	RootPath       string                 `json:"root_path"`
	PumaStatePaths []pumaStatusStatePaths `json:"puma_state_paths"`
}

type pumaStatusStatePaths struct {
	PumaStatePath       string                  `json:"puma_state_path"`
	BootedWorkers       int                     `json:"booted_workers"`
	AppCurrentPhase     int                     `json:"app_current_phase"`
	OldWorkers          int                     `json:"old_workers"`
	Workers             []pumaStatusWorker      `json:"workers"`
	TotalMaxThreads     int                     `json:"total_max_threads"`
	TotalCurrentThreads int                     `json:"total_current_threads"`
	MainPid             int                     `json:"main_pid"`
	Padding             *pumaStatusStatePadding `json:"-"`
}

type pumaStatusStatePadding struct {
	Pid      int
	CPU      int
	CPUTimes int
	Memory   int
	Uptime   int
}

type pumaStatusWorker struct {
	ID             int     `json:"id"`
	IsBooted       bool    `json:"is_booted"`
	Pid            int     `json:"pid"`
	LastCheckin    string  `json:"last_checkin"`
	CurrentThreads int     `json:"current_threads"`
	MaxThreads     int     `json:"max_threads"`
	CPUPercent     float64 `json:"cpu_percent"`
	Memory         float64 `json:"memory"`
	CurrentPhase   int     `json:"current_phase"`
	Uptime         int64   `json:"uptime"`
	CPUTimes       int     `json:"cpu_times"`
}
