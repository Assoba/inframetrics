/*
 * Copyright 2017  Assoba S.A.S.
 *
 * Licensed under the Apache License, Version 2.0 (the "License"); you may not use this file except in compliance with the License. You may obtain a copy of the License at
 *
 * http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software distributed under the License is distributed on an "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied. See the License for the specific language governing permissions and limitations under the License.
 *
 */

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

func loadIoStats() (*iostats, *iostats, *iostats) {
	statsLock.RLock()
	defer statsLock.RUnlock()
	ma1re := ma1
	ma5re := iostats{devices: make(map[string]deviceIO)}
	var ma5total = 0.0
	for _, i := range ma5 {
		ma5re.user += i.user
		ma5re.nice += i.nice
		ma5re.system += i.system
		ma5re.iowait += i.iowait
		ma5re.steal += i.steal
		ma5re.idle += i.idle
		for dev, io := range i.devices {
			oldVal, ok := ma5re.devices[dev]
			if ok {
				newVal := deviceIO{
					bytes_write: oldVal.bytes_write + io.bytes_write,
					bytes_read:  oldVal.bytes_read + io.bytes_read,
				}
				ma5re.devices[dev] = newVal
			} else {
				ma5re.devices[dev] = io
			}
		}
		ma5total += 1
	}

	ma5re.user = ma5re.user / ma5total
	ma5re.nice = ma5re.nice / ma5total
	ma5re.system = ma5re.system / ma5total
	ma5re.iowait = ma5re.iowait / ma5total
	ma5re.steal = ma5re.steal / ma5total
	ma5re.idle = ma5re.idle / ma5total
	for dev, io := range ma5re.devices {
		newVal := deviceIO{
			bytes_write: io.bytes_write / ma5total,
			bytes_read:  io.bytes_read / ma5total,
		}
		ma5re.devices[dev] = newVal
	}

	ma15re := iostats{devices: make(map[string]deviceIO)}
	var ma15total = 0.0
	for _, i := range ma15 {
		ma15re.user += i.user
		ma15re.nice += i.nice
		ma15re.system += i.system
		ma15re.iowait += i.iowait
		ma15re.steal += i.steal
		ma15re.idle += i.idle
		ma15total += 1
		for dev, io := range i.devices {
			oldVal, ok := ma15re.devices[dev]
			if ok {
				newVal := deviceIO{
					bytes_write: oldVal.bytes_write + io.bytes_write,
					bytes_read:  oldVal.bytes_read + io.bytes_read,
				}
				ma15re.devices[dev] = newVal
			} else {
				ma15re.devices[dev] = io
			}
		}
	}
	ma15re.user = ma15re.user / ma15total
	ma15re.nice = ma15re.nice / ma15total
	ma15re.system = ma15re.system / ma15total
	ma15re.iowait = ma15re.iowait / ma15total
	ma15re.steal = ma15re.steal / ma15total
	ma15re.idle = ma15re.idle / ma15total
	for dev, io := range ma15re.devices {
		newVal := deviceIO{
			bytes_write: io.bytes_write / ma15total,
			bytes_read:  io.bytes_read / ma15total,
		}
		ma15re.devices[dev] = newVal
	}

	return ma1re, &ma5re, &ma15re
}

func addIoStat(i *iostats) {
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

func RunIoStats() {
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
				var kbrs = deviceLine[2]
				kbr, err := strconv.ParseFloat(kbrs, 64)
				if err != nil {
					panic(err)
				}
				var br = kbr * 1024.0
				var kbws = deviceLine[3]
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
			addIoStat(&stats)
		}
	}
}
