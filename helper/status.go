package helper

import (
	"bytes"
	"encoding/json"
	"fmt"
	"os/exec"
	"time"
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
	fmt.Printf("Date: %s\n", time.Now().Format(time.RFC1123Z))
	fmt.Println()
}

func printApplicationsContext(pst pumaStatus) error {
	for _, key := range pst.WorkerStatus {
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

		fmt.Printf("* PID: %d\tBooted: %t\t\tRunner index: %d\n", key.Pid, key.Booted, key.Index)
		fmt.Printf("\t- Last status %s -\n", key.LastCheckin)
		fmt.Printf("  Running: %d\tPool capacity: %d\tMax threads: %d\n", key.LastStatus.Running, key.LastStatus.PoolCapacity, key.LastStatus.MaxThreads)
		fmt.Printf("  CPU: %s%%\tMemory: %s MiB\tTotal time: %ss\n", cpu, mem, ttime)

		fmt.Println()
	}

	return nil
}

func printApplicationGroups() error {
	fmt.Println("----------- Application groups -----------")

	for appname, key := range CfgFile.Applications {
		//output, err := exec.Command(key.Path+pumactlPath, "-S", key.Path+pumastatePath, "stats").Output()
		output, err := exec.Command("cat", "/go/src/github.com/dimelo/puma-helper/output.txt").Output()
		toutput := bytes.TrimLeft(output, "Command stats sent success")
		if err != nil {
			return err
		}

		var pst pumaStatus
		if err := json.Unmarshal(toutput, &pst); err != nil {
			return err
		}

		fmt.Printf("-> %s application (%s)\n", appname, key.State)
		fmt.Printf("  App root: %s\n", key.Path)
		fmt.Printf("  Booted workers: %d\n\n", pst.BootedWorkers)

		if err := printApplicationsContext(pst); err != nil {
			return err
		}

		fmt.Println()
		fmt.Println()
	}

	return nil
}
