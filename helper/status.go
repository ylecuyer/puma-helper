package helper

import (
	"fmt"
	"strings"

	log "github.com/Sirupsen/logrus"
)

const (
	appGroupName  string = "/current"
	pumactlPath   string = appGroupName + "/bin/pumactl"
	pumastatePath string = appGroupName + "/tmp/pids/puma.state"
)

// RunStatus run all status logical command
func RunStatus() error {
	rsd, err := retrieveStatusData()
	if err != nil {
		return err
	}

	if JSONOutput {
		return rsd.printAndBuildJSON()
	}

	printStatusGlobalInformations()

	return rsd.printStatusApps()
}

// retrieveStatusData fetch all data from Puma instances/workers
// and return it as a struct
func retrieveStatusData() (*pumaStatusFinalOutput, error) {
	appCount := 0
	apps := []pumaStatusApplication{}

	for appname, key := range CfgFile.Applications {
		if Filter != "" && !strings.Contains(appname, Filter) {
			continue
		}

		pcpath := key.Path + pumactlPath
		if key.PumactlPath != "" {
			pcpath = key.PumactlPath
		}

		pspath := key.Path + pumastatePath
		if key.PumastatePath != "" {
			pspath = key.PumastatePath
		}

		ps, err := readPumaStats(pcpath, pspath)
		if err != nil {
			log.Warn(fmt.Sprintf("[%s] configuration is invalid. Error: %v\n\n", appname, err))
			continue
		}

		workers := []pumaStatusWorker{}
		for _, v := range ps.WorkerStatus {
			pid := int32(v.Pid)

			cpu, err := getCPUFromPID(pid)
			if err != nil {
				return nil, err
			}

			mem, err := getMemoryFromPID(pid)
			if err != nil {
				return nil, err
			}

			ttime, err := getTotalTimeFromPID(pid)
			if err != nil {
				return nil, err
			}

			worker := pumaStatusWorker{
				ID:             v.Index,
				IsBooted:       v.Booted,
				Pid:            v.Pid,
				LastCheckin:    v.LastCheckin,
				CurrentThreads: v.LastStatus.MaxThreads - v.LastStatus.PoolCapacity,
				MaxThreads:     v.LastStatus.MaxThreads,
				CPUPercent:     cpu,
				Memory:         mem,
				TotalTimeExec:  int(ttime),
				CurrentPhase:   v.Phase,
			}

			workers = append(workers, worker)
		}

		app := pumaStatusApplication{
			Name:            appname,
			Description:     key.Description,
			RootPath:        key.Path,
			PumaStatePath:   pspath,
			PumaCtlPath:     pcpath,
			BootedWorkers:   ps.BootedWorkers,
			Worker:          workers,
			AppCurrentPhase: ps.Phase,
		}

		apps = append(apps, app)
		appCount++
	}

	return &pumaStatusFinalOutput{
		ApplicationsCount: appCount,
		Application:       apps,
	}, nil
}
