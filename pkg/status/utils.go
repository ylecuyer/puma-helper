package status

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net"
	"net/http"
	"strconv"
	"strings"
	"time"

	. "github.com/logrusorgru/aurora"
	proc "github.com/shirou/gopsutil/process"
	yaml "gopkg.in/yaml.v2"
)

const (
	bLoad string = "█"
	wLoad string = "░"
)

type pumaStateFile struct {
	Pid              int    `yaml:"pid"`
	ControlURL       string `yaml:"control_url"`
	ControlAuthToken string `yaml:"control_auth_token"`
}

// readPumaStats get and unmarshal JSON using puma unix socket
func readPumaStats(paths []string) ([]pumaStatus, error) {

	var pspaths []pumaStatus

	for _, path := range paths {
		data, err := ioutil.ReadFile(path)
		if err != nil {
			return nil, err
		}

		var psf pumaStateFile

		if err := yaml.Unmarshal(data, &psf); err != nil {
			return nil, err
		}

		httpc := http.Client{
			Transport: &http.Transport{
				Dial: func(_, _ string) (net.Conn, error) {
					return net.Dial("unix", strings.TrimPrefix(psf.ControlURL, "unix://"))
				},
			},
		}

		r, err := httpc.Get(fmt.Sprintf("http://unix/stats?token=%s", psf.ControlAuthToken))
		if err != nil {
			return nil, err
		}

		body, err := ioutil.ReadAll(r.Body)
		if err != nil {
			return nil, err
		}

		ps := pumaStatus{
			MainPid: psf.Pid,
		}

		if err := json.Unmarshal(body, &ps); err != nil {
			return nil, err
		}

		pspaths = append(pspaths, ps)
	}

	//body, _ := exec.Command("cat", "/go/src/github.com/dimelo/puma-helper/output.txt").Output()
	//fmt.Println(pspath)

	return pspaths, nil
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

func getCPUPercentFromPID(pid int32) (float64, error) {
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

func getCPUTimesFromPID(pid int32) (int, error) {
	p, err := proc.NewProcess(pid)
	if err != nil {
		return 0.0, err
	}
	t, err := p.Times()
	if err != nil {
		return 0.0, err
	}
	return strconv.Atoi(fmt.Sprintf("%.0f", t.Total()-(t.Idle+t.Iowait)))
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

	return elapsedFormatted(tx, time.Now().UTC())
}

func timeElapsedFromSeconds(t int) string {
	return timeElapsed(time.Now().UTC().Add(time.Duration(int64(t)) * time.Second).Format(time.RFC3339))
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

func elapsedFormatted(from, to time.Time) string {
	_, years, months, days, hours, minutes, seconds, _ := elapsed(from, to)

	estr := ""

	if years > 0 {
		estr = fmt.Sprintf("%dy%dm%dd", years, months, days)
	} else if months > 0 {
		estr = fmt.Sprintf("%dm%dd%dh", months, days, hours)
	} else if days > 0 {
		estr = fmt.Sprintf("%dd%dh%dm", days, hours, minutes)
	} else if hours > 0 {
		estr = fmt.Sprintf("%dh%dm%ds", hours, minutes, seconds)
	} else if minutes > 0 {
		estr = fmt.Sprintf("%dm%ds", minutes, seconds)
	} else {
		estr = fmt.Sprintf("%ds", seconds)
	}

	return estr
}

func elapsed(from, to time.Time) (inverted bool, years, months, days, hours, minutes, seconds, nanoseconds int) {
	if from.Location() != to.Location() {
		to = to.In(to.Location())
	}

	inverted = false
	if from.After(to) {
		inverted = true
		from, to = to, from
	}

	y1, M1, d1 := from.Date()
	y2, M2, d2 := to.Date()

	h1, m1, s1 := from.Clock()
	h2, m2, s2 := to.Clock()

	ns1, ns2 := from.Nanosecond(), to.Nanosecond()

	years = y2 - y1
	months = int(M2 - M1)
	days = d2 - d1

	hours = h2 - h1
	minutes = m2 - m1
	seconds = s2 - s1
	nanoseconds = ns2 - ns1

	if nanoseconds < 0 {
		nanoseconds += 1e9
		seconds--
	}
	if seconds < 0 {
		seconds += 60
		minutes--
	}
	if minutes < 0 {
		minutes += 60
		hours--
	}
	if hours < 0 {
		hours += 24
		days--
	}
	if days < 0 {
		days += daysIn(y2, M2-1)
		months--
	}
	if months < 0 {
		months += 12
		years--
	}
	return
}

func daysIn(year int, month time.Month) int {
	return time.Date(year, month+1, 0, 0, 0, 0, 0, time.UTC).Day()
}
