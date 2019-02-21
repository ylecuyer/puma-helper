package helper

import (
	"fmt"
	"strings"

	log "github.com/Sirupsen/logrus"
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

		pspath := key.Path + pumastatePath
		if key.PumastatePath != "" {
			pspath = key.PumastatePath
		}

		ps, err := readPumaStats(pspath)
		if err != nil {
			log.Warn(fmt.Sprintf("[%s] configuration is invalid. Error: %v\n\n", appname, err))
			continue
		}

		tmthreads := 0
		tcthreads := 0
		workers := []pumaStatusWorker{}
		for _, v := range ps.WorkerStatus {
			pid := int32(v.Pid)

			cpu, err := getCPUUsageFromPID(pid)
			if err != nil {
				return nil, err
			}

			mem, err := getMemoryFromPID(pid)
			if err != nil {
				return nil, err
			}

			ttime, err := getTotalExecTimeFromPID(pid)
			if err != nil {
				return nil, err
			}

			// Assuming this timestamp is in milliseconds
			utime, err := getTotalUptimeFromPID(pid)
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
				Uptime:         utime / 1000,
			}

			tcthreads += (v.LastStatus.MaxThreads - v.LastStatus.PoolCapacity)
			tmthreads += v.LastStatus.MaxThreads

			workers = append(workers, worker)
		}

		app := pumaStatusApplication{
			Name:                appname,
			Description:         key.Description,
			RootPath:            key.Path,
			PumaStatePath:       pspath,
			BootedWorkers:       ps.BootedWorkers,
			Worker:              workers,
			AppCurrentPhase:     ps.Phase,
			TotalCurrentThreads: tcthreads,
			TotalMaxThreads:     tmthreads,
			MainPid:             ps.MainPid,
		}

		apps = append(apps, app)
		appCount++
	}

	return &pumaStatusFinalOutput{
		ApplicationsCount: appCount,
		Application:       apps,
	}, nil
}
