package helper

import (
	"encoding/json"
	"fmt"
	"time"

	. "github.com/logrusorgru/aurora"
)

func printStatusGlobalInformations() {
	fmt.Println("---------- Global informations ----------")
	fmt.Printf("Version: %s\n", Version)
	fmt.Printf("Date: %s\n\n\n", time.Now().Format(time.RFC1123Z))
}

// printStatusApps print apps context one by one
func (ps pumaStatusFinalOutput) printStatusApps() error {
	fmt.Println("----------- Application groups -----------")

	line := 0
	for _, key := range ps.Application {
		currentApp = key.Name

		fmt.Printf("-> %s application\n", key.Name)
		if key.Description != "" {
			fmt.Printf("  About: %s\n", key.Description)
		}
		fmt.Printf("  App root: %s\n", key.RootPath)
		fmt.Printf("  Booted workers: %d\n", key.BootedWorkers)
		fmt.Printf("  Current phase: %d | Old workers: %d\n\n", key.AppCurrentPhase, key.OldWorkers)

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
			phase = Red(fmt.Sprintf("%d != %d app", key.CurrentPhase, currentPhase))
		}

		bootbtn := BgGreen(Bold("[UP]"))
		if !key.IsBooted {
			bootbtn = BgRed(Bold("[DOWN]"))
			fmt.Printf("*  %s ~ PID %d\tWorker ID %d\tLast checkin: %s\n", bootbtn, key.Pid, key.ID, timeElapsed(key.LastCheckin))
			continue
		}

		fmt.Printf("*  %s ~ PID %d\t\tWorker ID %d\tCPU: %s%%\tMem: %s MiB\tPhase: %s\n", bootbtn, key.Pid, key.ID, colorCPU(key.CPUPercent), colorMemory(key.Memory), phase)
		fmt.Printf("  Active threads: %s\tLast checkin: %s\tTotal exec time: %s\n", asciiThreadLoad(key.CurrentThreads, key.MaxThreads), timeElapsed(key.LastCheckin), timeElapsed(time.Now().Add(time.Duration(-int64(key.TotalTimeExec))*time.Second).Format(time.RFC3339)))
	}
	return nil
}

// printAndBuildJSON marshal and print pumaStatusFinalOutput
func (ps pumaStatusFinalOutput) printAndBuildJSON() error {
	b, err := json.Marshal(ps)
	if err != nil {
		return err
	}

	fmt.Println(string(b))

	return nil
}