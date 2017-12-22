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
	log "github.com/sirupsen/logrus"
	"github.com/hashicorp/go-sockaddr"
	"os/exec"
	"bufio"
	"strings"
	"strconv"
	"sync"
)

type netStats struct {
	devices map[string]netDevIO
}

type netDevIO struct {
	bps_in  float64
	bps_out float64
}

var nstatsLock = sync.RWMutex{}
var nma1 *netStats = nil
var nma5 []*netStats = make([]*netStats, 0, 6)
var nma15 []*netStats = make([]*netStats, 0, 16)

func loadNetStats() (*netStats, *netStats, *netStats) {
	nstatsLock.RLock()
	defer nstatsLock.RUnlock()
	ma1re := nma1
	ma5re := netStats{devices: make(map[string]netDevIO)}
	var ma5total = 0.0
	for _, i := range nma5 {
		for dev, io := range i.devices {
			oldVal, ok := ma5re.devices[dev]
			if ok {
				newVal := netDevIO{
					bps_in:oldVal.bps_in+io.bps_in,
					bps_out:oldVal.bps_out+io.bps_out,
				}
				ma5re.devices[dev] = newVal
			} else {
				ma5re.devices[dev] = io
			}
		}
		ma5total += 1
	}
	for dev, io := range ma5re.devices {
		newVal := netDevIO{
			bps_in: io.bps_in/ ma5total,
			bps_out:  io.bps_out/ ma5total,
		}
		ma5re.devices[dev] = newVal
	}
	ma15re := netStats{devices: make(map[string]netDevIO)}
	var ma15total = 0.0
	for _, i := range nma15 {
		for dev, io := range i.devices {
			oldVal, ok := ma15re.devices[dev]
			if ok {
				newVal := netDevIO{
					bps_in:oldVal.bps_in+io.bps_in,
					bps_out:oldVal.bps_out+io.bps_out,
				}
				ma15re.devices[dev] = newVal
			} else {
				ma15re.devices[dev] = io
			}
		}
		ma15total += 1
	}
	for dev, io := range ma15re.devices {
		newVal := netDevIO{
			bps_in: io.bps_in/ ma15total,
			bps_out:  io.bps_out/ ma15total,
		}
		ma15re.devices[dev] = newVal
	}
	return ma1re, &ma5re, &ma15re
}

func addNetStat(n *netStats) {
	nstatsLock.Lock()
	defer nstatsLock.Unlock()
	nma1 = n
	nma5 = append(nma5, n)
	if len(nma5) >= 6 {
		tmp := make([]*netStats, 5, 6)
		for idx := 1; idx < 6; idx++ {
			tmp[idx-1] = nma5[idx]
		}
		nma5 = tmp
	}
	nma15 = append(nma15, n)
	if len(nma15) >= 16 {
		tmp := make([]*netStats, 15, 16)
		for idx := 1; idx < 16; idx++ {
			tmp[idx-1] = nma15[idx]
		}
		nma15 = tmp
	}
}

func RunNetStats() {
	addrs, err := sockaddr.GetAllInterfaces()
	if err != nil {
		panic(err)
	}
	match, _ := sockaddr.FilterIfByType(addrs, sockaddr.TypeIPv4)
	_, match, err = sockaddr.IfByFlag(`loopback`, match)
	if err != nil {
		panic(err)
	}
	var ifaces = make([]string, 0)
	for _, i := range match {
		log.WithField("Address", i.Name).Info("Interface found")
		ifaces = append(ifaces, i.Name)
	}
	log.Debug("Starting ifstat")
	var cmd = exec.Command("/usr/bin/ifstat", "-n", "-q", "-i", strings.Join(ifaces, ","), "1")
	out, err := cmd.StdoutPipe()
	if err != nil {
		panic(err)
	}
	cmd.Start()
	log.WithField("command", cmd).WithField("output", out).Debug("Ifstat started")
	scanner := bufio.NewScanner(out)
	scanner.Scan()
	scanner.Scan()
	for scanner.Scan() {
		var vals = strings.Fields(scanner.Text())
		log.WithField("Fields", vals).Debug("Got line")
		if len(vals) == len(ifaces)*2 {
			s := netStats{devices: make(map[string]netDevIO)}
			for i, a := range ifaces {
				in, err := strconv.ParseFloat(vals[i*2+0], 64)
				if err != nil {
					panic(err)
				}
				out, err := strconv.ParseFloat(vals[i*2+1], 64)
				if err != nil {
					panic(err)
				}
				d := netDevIO{
					bps_in:  in * 1024,
					bps_out: out * 1024,
				}
				s.devices[a] = d
			}
			log.WithField("Stats", s).Debug("Collected netStats")
			addNetStat(&s)
		}
	}
}
