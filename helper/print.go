package helper

import (
	"encoding/json"
	"fmt"
	"time"

	. "github.com/logrusorgru/aurora"
)

func printStatusGlobalInformations() {
	fmt.Println("---------- Global informations ----------")
	if ExpandDetails {
		fmt.Printf("Version: %s\n", Version)
		fmt.Printf("Date: %s\n\n\n", time.Now().Format(time.RFC1123Z))
		return
	}
	fmt.Printf("Version: %s | Date: %s\n\n", Version, time.Now().Format(time.RFC1123Z))
}

// printStatusApps print apps context one by one
func (ps pumaStatusFinalOutput) printStatusApps() error {
	fmt.Println("----------- Application groups -----------")

	line := 0
	for _, key := range ps.Application {
		currentApp = key.Name

		if ExpandDetails {
			fmt.Printf("-> %s application\n", key.Name)
			if key.Description != "" {
				fmt.Printf("  About: %s\n", key.Description)
			}
			fmt.Printf("  App root: %s\n", key.RootPath)
			fmt.Printf("  Booted workers: %d | PID: %d\n", key.BootedWorkers, key.MainPid)
			fmt.Printf("  Current phase: %d | Old workers: %d | Active threads: %s\n\n", key.AppCurrentPhase, key.OldWorkers, asciiThreadLoad(key.TotalCurrentThreads, key.TotalMaxThreads))
		} else {
			if key.OldWorkers > 0 {
				fmt.Printf("-> %s %d (%s) Phase: %d | Workers: %d (Old: %d) | Active threads: %s\n\n", key.Name, key.MainPid, key.RootPath, key.AppCurrentPhase, key.BootedWorkers, key.OldWorkers, asciiThreadLoad(key.TotalCurrentThreads, key.TotalMaxThreads))
			} else {
				fmt.Printf("-> %s %d (%s) Phase: %d | Workers: %d | Active threads: %s\n\n", key.Name, key.MainPid, key.RootPath, key.AppCurrentPhase, key.BootedWorkers, asciiThreadLoad(key.TotalCurrentThreads, key.TotalMaxThreads))
			}
		}

		if err := printStatusWorkers(key.Worker, key.AppCurrentPhase); err != nil {
			return err
		}

		if line < len(ps.Application)-1 {
			fmt.Println()
			line++
		}
	}

	return nil
}

// printStatusWorkers print workers status context of one app
func printStatusWorkers(ps []pumaStatusWorker, currentPhase int) error {
	for _, key := range ps {
		phase := Green(fmt.Sprintf("%d", key.CurrentPhase))
		if key.CurrentPhase != currentPhase {
			phase = Red(fmt.Sprintf("%d", key.CurrentPhase))
		}

		if !ExpandDetails {
			lcheckin := Green(timeElapsed(key.LastCheckin))
			if len(timeElapsed(key.LastCheckin)) >= 3 {
				lcheckin = Brown(timeElapsed(key.LastCheckin))
			}

			fmt.Printf("* %d [%d] CPU Av: %s%% CPU Times: %s Mem: %sMiB Phase: %s Uptime: %s Threads: %s (Last checkin: %s)\n", key.ID, key.Pid, colorCPU(key.CPUPercent), timeElapsedFromSeconds(key.CPUTimes), colorMemory(key.Memory), phase, timeElapsed(time.Unix(key.Uptime, 0).Format(time.RFC3339)), asciiThreadLoad(key.CurrentThreads, key.MaxThreads), lcheckin)
			continue
		}

		bootbtn := BgGreen(Bold("[UP]"))
		if !key.IsBooted {
			bootbtn = BgRed(Bold("[DOWN]"))
			fmt.Printf("*  %s ~ PID %d\tWorker ID %d\tLast checkin: %s\n", bootbtn, key.Pid, key.ID, timeElapsed(key.LastCheckin))
			continue
		}

		fmt.Printf("*  %s ~ PID %d\tWorker ID %d\tCPU Average: %s%%\tMem: %sMiB\tActive threads: %s\n", bootbtn, key.Pid, key.ID, colorCPU(key.CPUPercent), colorMemory(key.Memory), asciiThreadLoad(key.CurrentThreads, key.MaxThreads))
		fmt.Printf("  Phase: %s\tLast checkin: %s\tTotal CPU times: %s\tUptime: %s\n", phase, timeElapsed(key.LastCheckin), timeElapsedFromSeconds(key.CPUTimes), timeElapsed(time.Unix(key.Uptime, 0).Format(time.RFC3339)))
	}
	return nil
}

// printAndBuildJSON marshal and print pumaStatusFinalOutput struct
func (ps pumaStatusFinalOutput) printAndBuildJSON() error {
	b, err := json.Marshal(ps)
	if err != nil {
		return err
	}

	fmt.Println(string(b))

	return nil
}
