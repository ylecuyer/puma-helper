package helper

import (
	"encoding/json"
	"fmt"
	"log"
	"os/exec"
	"time"
)

const (
	appGroupName  string = "/current"
	pumactlPath   string = appGroupName + "/bin/pumactl"
	pumastatePath string = appGroupName + "/tmp/pids/puma.state"
)

var (
	ps pumaStatus
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

// RunStatus ...
func RunStatus() {
	printGlobalInformations()
	printApplicationGroups()
}

func printGlobalInformations() {
	fmt.Println("---------- Global informations ----------")
	fmt.Printf("Version: %s\n", Version)
	fmt.Printf("Date: %s\n", time.Now().Format(time.RFC822Z))
	fmt.Println()
	fmt.Println(getMemoryFromPID(696))
}

func printApplicationGroups() {
	fmt.Println("----------- Application groups -----------")

	for appname, k := range CfgFile.Applications {
		fmt.Printf("-> %s application (%s)\n", appname, k.State)
		fmt.Printf("  App root: %s\n\n", k.Path)

		output, err := exec.Command(k.Path+pumactlPath, "-S", k.Path+pumastatePath, "stats").Output()
		if err != nil {
			log.Fatal(err)
		}
		json.Unmarshal(output, &ps)

		for _, m := range ps.WorkerStatus {
			fmt.Printf("* PID: %d\tUptime: %s\n", m.Pid, m.LastCheckin)
			fmt.Printf("  CPU: 0P\tMemory: 231M\n")
			fmt.Println()
		}
		fmt.Println()
		fmt.Println()
		fmt.Println()
		fmt.Println()
		fmt.Println()
	}
}
