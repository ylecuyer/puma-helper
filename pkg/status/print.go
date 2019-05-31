package status

import (
	"encoding/json"
	"fmt"
	"strconv"
	"strings"
	"time"

	version "github.com/dimelo/puma-helper/pkg/version"
	. "github.com/logrusorgru/aurora"
	log "github.com/sirupsen/logrus"
)

func printStatusHeader() {
	fmt.Printf("Version: %s | Date: %s\n\n", version.Version, time.Now().Format(time.RFC1123Z))
}

// printStatusApps print apps context one by one
func (ps pumaStatusFinalOutput) printStatusApps() {
	line := 0

	for _, key := range ps.Application {
		currentApp = key.Name

		basetitle := fmt.Sprintf("-> %s application", currentApp)
		if ExpandDetails {
			fmt.Printf(basetitle + " | App root: " + key.RootPath)
		}
		if len(key.Description) == 0 {
			fmt.Printf(basetitle)
		} else {
			fmt.Printf(basetitle + " | About: " + key.Description)
		}

		for _, keypath := range key.PumaStatePaths {
			if ExpandDetails {
				fmt.Printf("\n  -> File: %s\n", keypath.PumaStatePath)

				fmt.Printf("  Booted workers: %d | PID: %d\n", keypath.BootedWorkers,
					keypath.MainPid)

				fmt.Printf("  Current phase: %d | Old workers: %d | Load: %s\n\n",
					keypath.AppCurrentPhase,
					keypath.OldWorkers,
					asciiThreadLoad(keypath.TotalCurrentThreads, keypath.TotalMaxThreads))
			} else {
				fmt.Printf("\n  %d (%s) | Load: %s\n",
					keypath.MainPid,
					keypath.PumaStatePath,
					asciiThreadLoad(keypath.TotalCurrentThreads, keypath.TotalMaxThreads))
			}

			printStatusWorkers(keypath, keypath.AppCurrentPhase)
		}

		if line < len(ps.Application)-1 {
			fmt.Println()
			line++
		}
	}
}

// printStatusWorkers print workers status context of one app
func printStatusWorkers(pstatuspath pumaStatusStatePaths, currentPhase int) {

	ps := pstatuspath.Workers
	pad := *pstatuspath.Padding

	padcpu := strconv.Itoa(pad.CPU)
	padcput := strconv.Itoa(pad.CPUTimes)
	padmem := strconv.Itoa(pad.Memory)
	paduptime := strconv.Itoa(pad.Uptime)
	padpid := strconv.Itoa(pad.Pid)

	for _, key := range ps {
		phase := Green(fmt.Sprintf("%d", key.CurrentPhase))
		if key.CurrentPhase != currentPhase {
			phase = Red(fmt.Sprintf("%d", key.CurrentPhase))
		}

		te := timeElapsed(key.LastCheckin)

		if !ExpandDetails {
			fmt.Printf("└ %"+padpid+"d %"+padcpu+"s%% %"+padmem+"sM Uptime: %"+paduptime+"s Load: %s",
				key.Pid, colorCPU(key.CPUPercent),
				colorMemory(key.Memory),
				timeElapsed(time.Unix(key.Uptime, 0).Format(time.RFC3339)),
				asciiThreadLoad(key.CurrentThreads, key.MaxThreads))

			if len(te) >= 3 || !strings.Contains(te, "s") {
				fmt.Printf(" %s", Red("Last checkin: "+te))
			}

			if key.CurrentPhase != currentPhase {
				fmt.Printf(" %s", Red("Phase: "+string(key.CurrentPhase)))
			}

			fmt.Println()

		} else {
			bootbtn := BgGreen(Bold("[UP]"))
			if !key.IsBooted {
				bootbtn = BgRed(Bold("[DOWN]"))
				fmt.Printf("*  %s ~ PID %d\tWorker ID %d\tLast checkin: %s\n", bootbtn, key.Pid, key.ID, te)
				continue
			}

			fmt.Printf("*  %s ~ PID %"+padpid+"d\tWorker ID %d\tCPU Average: %"+padcpu+"s%%\tMem: %"+padmem+"sM\tLoad: %s\n",
				bootbtn,
				key.Pid,
				key.ID,
				colorCPU(key.CPUPercent),
				colorMemory(key.Memory),
				asciiThreadLoad(key.CurrentThreads, key.MaxThreads))

			fmt.Printf("  Phase: %s\tLast checkin: %s\tTotal CPU times: %"+padcput+"s\tUptime: %"+paduptime+"s\n",
				phase,
				te,
				timeElapsedFromSeconds(key.CPUTimes),
				timeElapsed(time.Unix(key.Uptime, 0).Format(time.RFC3339)))
		}
	}
}

func printRetrieveStatusError() {
	for k, v := range retrieveStatusError {
		if k == 0 {
			fmt.Println()
		}

		log.Warn(v)
	}
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
