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
		fmt.Printf("  Booted workers: %d\n\n", key.BootedWorkers)

		if err := printStatusWorkers(key.Worker); err != nil {
			return err
		}

		if line < len(ps.Application)-1 {
			fmt.Println()
			line++
		}
	}

	return nil
}

func printStatusWorkers(ps []pumaStatusWorker) error {
	for _, key := range ps {
		bootbtn := BgGreen(Bold("[UP]"))
		if !key.IsBooted {
			bootbtn = BgRed(Bold("[DOWN]"))
			fmt.Printf("*  %s ~ PID %d\tWorker ID %d\tLast checkin: %s\n", bootbtn, key.Pid, key.ID, timeElapsed(key.LastCheckin))
			continue
		}

		fmt.Printf("*  %s ~ PID %d\t\tWorker ID %d\tCPU: %s%%\tMemory: %s MiB\n", bootbtn, key.Pid, key.ID, colorCPU(key.CPUPercent), colorMemory(key.Memory))
		fmt.Printf("  Active threads: %s\tLast checkin: %s\tTotal exec time: %s\n", asciiThreadLoad(key.CurrentThreads, key.MaxThreads), timeElapsed(key.LastCheckin), timeElapsed(time.Now().Add(time.Duration(-int64(key.TotalTimeExec))*time.Second).Format(time.RFC3339)))
	}
	return nil
}

func (ps pumaStatusFinalOutput) printAndBuildJSON() error {
	b, err := json.Marshal(ps)
	if err != nil {
		return err
	}

	fmt.Println(string(b))

	return nil
}
