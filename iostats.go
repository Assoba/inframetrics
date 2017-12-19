package inframetrics

import (
	"os/exec"
	"bufio"
	"strings"
	"strconv"
	"sync"
	log "github.com/sirupsen/logrus"
)

type deviceIO struct {
	tps         float64
	bytes_write float64
	bytes_read  float64
}

type iostats struct {
	user    float64
	nice    float64
	system  float64
	iowait  float64
	steal   float64
	idle    float64
	devices map[string]deviceIO
}

var statsLock = sync.RWMutex{}
var ma1 *iostats = nil
var ma5 []*iostats = make([]*iostats, 0, 6)
var ma15 []*iostats = make([]*iostats, 0, 16)

func GetStats() (*iostats, *iostats, *iostats) {
	statsLock.RLock()
	defer statsLock.RUnlock()
	ma1re := ma1
	ma5re := iostats{}
	var ma5total = 0.0
	for _, i := range ma5 {
		ma5re.user += i.user
		ma5re.nice += i.nice
		ma5re.system += i.system
		ma5re.iowait += i.iowait
		ma5re.steal += i.steal
		ma5re.idle += i.idle
		ma5total += 1
	}
	ma5re.user = ma5re.user / ma5total
	ma5re.nice = ma5re.nice / ma5total
	ma5re.system = ma5re.system / ma5total
	ma5re.iowait = ma5re.iowait / ma5total
	ma5re.steal = ma5re.steal / ma5total
	ma5re.idle = ma5re.idle / ma5total
	ma15re := iostats{}
	var ma15total = 0.0
	for _, i := range ma15 {
		ma15re.user += i.user
		ma15re.nice += i.nice
		ma15re.system += i.system
		ma15re.iowait += i.iowait
		ma15re.steal += i.steal
		ma15re.idle += i.idle
		ma15total += 1
	}
	ma15re.user = ma5re.user / ma15total
	ma15re.nice = ma5re.nice / ma15total
	ma15re.system = ma5re.system / ma15total
	ma15re.iowait = ma5re.iowait / ma15total
	ma15re.steal = ma5re.steal / ma15total
	ma15re.idle = ma5re.idle / ma15total
	return ma1re, &ma5re, &ma15re
}

func AddStat(i *iostats) {
	statsLock.Lock()
	defer statsLock.Unlock()
	ma1 = i
	ma5 = append(ma5, i)
	if len(ma5) >= 6 {
		tmp := make([]*iostats, 5, 6)
		for idx := 1; idx < 6; idx++ {
			tmp[idx-1] = ma5[idx]
		}
		ma5 = tmp
	}
	ma15 = append(ma15, i)
	if len(ma15) >= 16 {
		tmp := make([]*iostats, 15, 16)
		for idx := 1; idx < 16; idx++ {
			tmp[idx-1] = ma15[idx]
		}
		ma15 = tmp
	}
}

func RunStats() {
	log.Debug("Starting iostat")
	var cmd = exec.Command("/usr/bin/iostat", "-y", "1")
	out, err := cmd.StdoutPipe()
	if err != nil {
		panic(err)
	}
	cmd.Start()
	log.WithField("command", cmd).WithField("output", out).Debug("Iostat started")
	scanner := bufio.NewScanner(out)
	scanner.Scan()
	for scanner.Scan() {
		cpuLine := scanner.Text()
		var cpu = strings.Fields(cpuLine)
		if len(cpu) == 6 {
			log.WithField("cpuLine", cpu).Debug("Got CPU line")
			user, err := strconv.ParseFloat(cpu[0], 64)
			if err != nil {
				panic(err)
			}
			log.WithField("user", user).Debug("Parsed user")
			nice, err := strconv.ParseFloat(cpu[1], 64)
			if err != nil {
				panic(err)
			}
			log.WithField("nice", user).Debug("Parsed nice")
			system, err := strconv.ParseFloat(cpu[2], 64)
			if err != nil {
				panic(err)
			}
			iowait, err := strconv.ParseFloat(cpu[3], 64)
			if err != nil {
				panic(err)
			}
			log.WithField("iowait", user).Debug("Parsed iowait")
			steal, err := strconv.ParseFloat(cpu[4], 64)
			if err != nil {
				panic(err)
			}
			log.WithField("steal", user).Debug("Parsed steal")
			idle, err := strconv.ParseFloat(cpu[5], 64)
			if err != nil {
				panic(err)
			}
			log.WithField("idle", user).Debug("Parsed idle")
			var stats = iostats{
				user:    user,
				nice:    nice,
				system:  system,
				iowait:  iowait,
				steal:   steal,
				idle:    idle,
				devices: make(map[string]deviceIO),
			}
			log.WithField("stat", stats).Debug("Parsed cpu")
			// Skip next 2 lines
			scanner.Scan()
			scanner.Scan()
			// Read first device line
			scanner.Scan()
			var deviceLine = strings.Fields(scanner.Text())
			log.WithField("deviceLine", deviceLine).Debug("First device line")
			for len(deviceLine) == 6 {
				var dev = deviceLine[0]
				var tpss = deviceLine[1]
				tps, err := strconv.ParseFloat(tpss, 64)
				if err != nil {
					panic(err)
				}
				var kbrs = cpu[2]
				kbr, err := strconv.ParseFloat(kbrs, 64)
				if err != nil {
					panic(err)
				}
				var br = kbr * 1024.0
				var kbws = cpu[2]
				kbw, err := strconv.ParseFloat(kbws, 64)
				if err != nil {
					panic(err)
				}
				var bw = kbw * 1024.0
				stats.devices[dev] = deviceIO{
					tps:         tps,
					bytes_read:  br,
					bytes_write: bw,
				}
				scanner.Scan()
				deviceLine = strings.Fields(scanner.Text())
				log.WithField("deviceLine", deviceLine).Debug("Next device line")
			}
			AddStat(&stats)
		}
	}
}
