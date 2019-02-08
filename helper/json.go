package helper

type pumaStatus struct {
	BootedWorkers int `json:"booted_workers"`
	OldWorkers    int `json:"old_workers"`
	Phase         int `json:"phase"`
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
	Name            string             `json:"name"`
	Description     string             `json:"description"`
	RootPath        string             `json:"root_path"`
	PumaStatePath   string             `json:"puma_state_path"`
	PumaCtlPath     string             `json:"pumactl_path"`
	BootedWorkers   int                `json:"booted_workers"`
	AppCurrentPhase int                `json:"app_current_phase"`
	OldWorkers      int                `json:"old_workers"`
	Worker          []pumaStatusWorker `json:"worker"`
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
	TotalTimeExec  int     `json:"total_time_exec"`
	CurrentPhase   int     `json:"current_phase"`
	Uptime         int64   `json:"uptime"`
}
