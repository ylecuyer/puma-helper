package helper

import (
	"bytes"
	"encoding/json"
	"fmt"
	"os/exec"
	"strings"
	"time"

	. "github.com/logrusorgru/aurora"
	proc "github.com/shirou/gopsutil/process"
)

const (
	bLoad string = "█"
	wLoad string = "░"
)

// readPumaStats get and unmarshal JSON using pumactl
func readPumaStats(pcpath, pspath string) (pumaStatus, error) {
	var ps pumaStatus

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

func getTotalExecTimeFromPID(pid int32) (float64, error) {
	p, err := proc.NewProcess(pid)
	if err != nil {
		return 0.0, err
	}
	t, _ := p.Times()
	return t.Total(), nil
}

func getTotalUptimeFromPID(pid int32) (int64, error) {
	p, err := proc.NewProcess(pid)
	if err != nil {
		return 0, err
	}
	t, err := p.CreateTime()
	if err != nil {
		return 0, err
	}
	return t, nil
}

func getCPUFromPID(pid int32) (float64, error) {
	p, err := proc.NewProcess(pid)
	if err != nil {
		return 0.0, err
	}
	cpu, err := p.CPUPercent()
	if err != nil {
		return 0.0, err
	}
	return cpu, nil
}

func getMemoryFromPID(pid int32) (float64, error) {
	p, err := proc.NewProcess(pid)
	if err != nil {
		return 0.0, err
	}
	mem, err := p.MemoryInfoEx()
	if err != nil {
		return 0.0, err
	}
	return float64(mem.RSS+mem.Shared) / float64(1024*1024), nil
}

func timeElapsed(nT string) string {
	tx, err := time.Parse(time.RFC3339, nT)
	if err != nil {
		return "unrecognized time format"
	}

	elapsed := time.Since(tx).String()
	if strings.Contains(elapsed, "ms") {
		return "~0s"
	}

	return fmt.Sprintf("%ss", strings.Split(elapsed, ".")[0])
}

func colorState(fvalue, warnstate, criticalstate float64, strvalue string) string {
	if fvalue > criticalstate {
		return Red(strvalue).String()
	} else if fvalue > warnstate {
		return Brown(strvalue).String()
	}

	return Green(strvalue).String()
}

func asciiThreadLoad(load int, capacity int) string {
	formatted := fmt.Sprintf("%d[%s%s]%d", load, strings.Repeat(bLoad, load), strings.Repeat(wLoad, capacity-load), capacity)
	total := (float64(load) / float64(capacity)) * 100

	pwarn := CfgFile.Applications[currentApp].ThreadWarn
	if pwarn == 0 {
		pwarn = 50
	}
	pcritical := CfgFile.Applications[currentApp].ThreadCritical
	if pcritical == 0 {
		pcritical = 80
	}

	return colorState(total, float64(pwarn), float64(pcritical), formatted)
}

func colorCPU(cpu float64) string {
	cwarn := CfgFile.Applications[currentApp].CPUWarn
	if cwarn == 0 {
		cwarn = 50
	}
	ccritical := CfgFile.Applications[currentApp].CPUCritical
	if ccritical == 0 {
		ccritical = 80
	}

	return colorState(cpu, float64(cwarn), float64(ccritical), fmt.Sprintf("%.1f", cpu))
}

func colorMemory(memory float64) string {
	mwarn := CfgFile.Applications[currentApp].MemoryWarn
	if mwarn == 0 {
		mwarn = 500
	}
	mcritical := CfgFile.Applications[currentApp].MemoryCritical
	if mcritical == 0 {
		mcritical = 1000
	}

	return colorState(memory, float64(mwarn), float64(mcritical), fmt.Sprintf("%.1f", memory))
}
