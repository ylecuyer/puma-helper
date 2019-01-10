package helper

import (
	"bytes"
	"encoding/json"
	"fmt"
	"os/exec"
	"strings"
	"time"

	. "github.com/logrusorgru/aurora"
)

const (
	appGroupName  string = "/current"
	pumactlPath   string = appGroupName + "/bin/pumactl"
	pumastatePath string = appGroupName + "/tmp/pids/puma.state"
)

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

// RunStatus run all status logical command
func RunStatus() error {
	printGlobalInformations()
	if err := printApplicationGroups(); err != nil {
		return err
	}
	return nil
}

func printGlobalInformations() {
	fmt.Println("---------- Global informations ----------")
	fmt.Printf("Version: %s\n", Version)
	fmt.Printf("Date: %s\n\n\n", time.Now().Format(time.RFC1123Z))
}

func printApplicationsContext(pst pumaStatus) error {
	for _, key := range pst.WorkerStatus {
		bootbtn := BgGreen(Bold("[UP]"))
		if !key.Booted {
			bootbtn = BgRed(Bold("[DOWN]"))
			fmt.Printf("*  %s ~ PID %d / Worker ID %d\tLast checkin: %s\n", bootbtn, key.Pid, key.Index, timeElapsed(key.LastCheckin))
			continue
		}

		cpu, err := getCPUFromPID(int32(key.Pid))
		if err != nil {
			return err
		}

		mem, err := getMemoryFromPID(int32(key.Pid))
		if err != nil {
			return err
		}

		ttime, err := getTotalTimeFromPID(int32(key.Pid))
		if err != nil {
			return err
		}

		fmt.Printf("*  %s ~ PID %d / Worker ID %d\tActive threads: [%d / %d]\tLast checkin: %s\n", bootbtn, key.Pid, key.Index, key.LastStatus.Running, key.LastStatus.PoolCapacity, timeElapsed(key.LastCheckin))
		fmt.Printf("  CPU: %s%%\tMemory: %s MiB\tTotal exec time: %s\n", cpu, mem, timeElapsed(time.Now().Add(time.Duration(-int64(ttime))*time.Second).Format(time.RFC3339)))
	}

	return nil
}

func readPumaStats(key Application) (pumaStatus, error) {
	var ps pumaStatus

	pcpath := key.Path + pumactlPath
	if key.PumactlPath != "" {
		pcpath = key.PumactlPath
	}

	pspath := key.Path + pumastatePath
	if key.PumastatePath != "" {
		pspath = key.PumastatePath
	}

	output, err := exec.Command(pcpath, "-S", pspath, "stats").Output()
	//fmt.Println(pcpath, pspath)
	//output, err := exec.Command("cat", "/go/src/github.com/dimelo/puma-helper/output.txt").Output()
	if err != nil {
		return ps, err
	}

	toutput := bytes.TrimLeft(output, "Command stats sent success")

	if err := json.Unmarshal(toutput, &ps); err != nil {
		return ps, err
	}

	return ps, nil
}

func printApplicationGroups() error {
	fmt.Println("----------- Application groups -----------")

	line := 0
	for appname, key := range CfgFile.Applications {
		if Filter != "" && !strings.Contains(appname, Filter) {
			continue
		}

		ps, err := readPumaStats(key)
		if err != nil {
			return err
		}

		fmt.Printf("-> %s application\n", appname)
		if key.Description != "" {
			fmt.Printf("  About: %s\n", key.Description)
		}
		fmt.Printf("  App root: %s\n", key.Path)
		fmt.Printf("  Booted workers: %d\n\n", ps.BootedWorkers)

		if err := printApplicationsContext(ps); err != nil {
			return err
		}

		if line < len(CfgFile.Applications)-1 {
			fmt.Println()
			line++
		}
	}

	return nil
}
