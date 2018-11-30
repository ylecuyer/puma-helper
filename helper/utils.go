package helper

import (
	"fmt"

	proc "github.com/shirou/gopsutil/process"
)

func getTotalTimeFromPID(pid int32) (string, error) {
	p, err := proc.NewProcess(pid)
	if err != nil {
		return "", err
	}
	t, _ := p.Times()
	return fmt.Sprintf("%.2f", t.Total()), nil
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
	return fmt.Sprintf("%.2f", cpu), nil
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
	return fmt.Sprintf("%.2f", float64(mem.RSS+mem.Shared)/float64(1024*1024)), nil
}
