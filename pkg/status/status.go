package status

import (
	"fmt"
	"sort"
	"strings"

	log "github.com/sirupsen/logrus"
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

	sortedApps := make([]string, 0, len(CfgFile.Applications))
	for k := range CfgFile.Applications {
		sortedApps = append(sortedApps, k)
	}
	sort.Strings(sortedApps)

	for _, appName := range sortedApps {
		key := CfgFile.Applications[appName]

		if Filter != "" && !strings.Contains(appName, Filter) {
			continue
		}

		var pspath []string
		pspath = append(pspath, key.Path+pumastatePath)
		if len(key.PumastatePaths) > 0 {
			pspath = key.PumastatePaths
		}

		pss, err := readPumaStats(pspath)
		if err != nil {
			log.Warn(fmt.Sprintf("[%s] configuration is invalid. Error: %v\n\n", appName, err))
			continue
		}

		pssps := []pumaStatusStatePaths{}

		for fid, ps := range pss {

			workers := []pumaStatusWorker{}
			tmthreads := 0
			tcthreads := 0

			for _, v := range ps.WorkerStatus {
				pid := int32(v.Pid)

				cpu, err := getCPUPercentFromPID(pid)
				if err != nil {
					return nil, err
				}

				cput, err := getCPUTimesFromPID(pid)
				if err != nil {
					return nil, err
				}

				mem, err := getMemoryFromPID(pid)
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
					CurrentPhase:   v.Phase,
					Uptime:         utime / 1000,
					CPUTimes:       cput,
				}

				tcthreads += (v.LastStatus.MaxThreads - v.LastStatus.PoolCapacity)
				tmthreads += v.LastStatus.MaxThreads

				workers = append(workers, worker)
			}

			pssp := pumaStatusStatePaths{
				PumaStatePath:       pspath[fid],
				BootedWorkers:       ps.BootedWorkers,
				Workers:             workers,
				TotalCurrentThreads: tcthreads,
				TotalMaxThreads:     tmthreads,
				MainPid:             ps.MainPid,
				AppCurrentPhase:     ps.Phase,
			}

			pssps = append(pssps, pssp)
		}

		app := pumaStatusApplication{
			Name:           appName,
			Description:    key.Description,
			RootPath:       key.Path,
			PumaStatePaths: pssps,
		}

		apps = append(apps, app)
		appCount++
	}

	return &pumaStatusFinalOutput{
		ApplicationsCount: appCount,
		Application:       apps,
	}, nil
}
