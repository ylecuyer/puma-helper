package helper

import (
	"fmt"
	"strings"
	"time"

	proc "github.com/shirou/gopsutil/process"
)

func getTotalTimeFromPID(pid int32) (float64, error) {
	p, err := proc.NewProcess(pid)
	if err != nil {
		return 0.0, err
	}
	t, _ := p.Times()
	return t.Total(), nil
}

func getCPUFromPID(pid int32) (string, error) {
	p, err := proc.NewProcess(pid)
	if err != nil {
		return "", err
	}
	cpu, err := p.CPUPercent()
	if err != nil {
		return "", err
	}
	return fmt.Sprintf("%.1f", cpu), nil
}

func getMemoryFromPID(pid int32) (string, error) {
	p, err := proc.NewProcess(pid)
	if err != nil {
		return "", err
	}
	mem, err := p.MemoryInfoEx()
	if err != nil {
		return "", err
	}
	return fmt.Sprintf("%.1f", float64(mem.RSS+mem.Shared)/float64(1024*1024)), nil
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
