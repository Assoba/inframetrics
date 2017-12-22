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
	"net/http"
	"bufio"
	"strings"
	"github.com/shirou/gopsutil/mem"
	"strconv"
	"github.com/shirou/gopsutil/load"
	"os"
)

func Handler(w http.ResponseWriter, _ *http.Request) {
	getSystemMetrics(w)
	getIoStats(w)
	getNetStats(w)
	getNomadMetrics(w)
	getConsulMetrics(w)
}

func getNetStats(w http.ResponseWriter) {
	m1, m5, m15 := loadNetStats()
	var crlf = []byte("\r\n")
	for dev, io := range m1.devices {
		w.Write([]byte("node_net_in_bytes_per_second{duration=\"1s\", device=\"" + dev + "\"} "))
		w.Write([]byte(strconv.FormatFloat(io.bps_in, 'g', 3, 64)))
		w.Write(crlf)
		w.Write([]byte("node_net_out_bytes_per_second{duration=\"1s\", device=\"" + dev + "\"} "))
		w.Write([]byte(strconv.FormatFloat(io.bps_out, 'g', 3, 64)))
		w.Write(crlf)
	}
	for dev, io := range m5.devices {
		w.Write([]byte("node_net_in_bytes_per_second{duration=\"5s\", device=\"" + dev + "\"} "))
		w.Write([]byte(strconv.FormatFloat(io.bps_in, 'g', 3, 64)))
		w.Write(crlf)
		w.Write([]byte("node_net_out_bytes_per_second{duration=\"5s\", device=\"" + dev + "\"} "))
		w.Write([]byte(strconv.FormatFloat(io.bps_out, 'g', 3, 64)))
		w.Write(crlf)
	}
	for dev, io := range m15.devices {
		w.Write([]byte("node_net_in_bytes_per_second{duration=\"15s\", device=\"" + dev + "\"} "))
		w.Write([]byte(strconv.FormatFloat(io.bps_in, 'g', 3, 64)))
		w.Write(crlf)
		w.Write([]byte("node_net_out_bytes_per_second{duration=\"15s\", device=\"" + dev + "\"} "))
		w.Write([]byte(strconv.FormatFloat(io.bps_out, 'g', 3, 64)))
		w.Write(crlf)
	}
}

func getIoStats(w http.ResponseWriter) {
	m1, m5, m15 := loadIoStats()
	var crlf = []byte("\r\n")
	w.Write([]byte("node_cpu_user_ratio{duration=\"1s\"} "))
	w.Write([]byte(strconv.FormatFloat(m1.user, 'g', 3, 64)))
	w.Write(crlf)
	w.Write([]byte("node_cpu_user_ratio{duration=\"5s\"} "))
	w.Write([]byte(strconv.FormatFloat(m5.user, 'g', 3, 64)))
	w.Write(crlf)
	w.Write([]byte("node_cpu_user_ratio{duration=\"15s\"} "))
	w.Write([]byte(strconv.FormatFloat(m15.user, 'g', 3, 64)))
	w.Write(crlf)

	w.Write([]byte("node_cpu_nice_ratio{duration=\"1s\"} "))
	w.Write([]byte(strconv.FormatFloat(m1.nice, 'g', 3, 64)))
	w.Write(crlf)
	w.Write([]byte("node_cpu_nice_ratio{duration=\"5s\"} "))
	w.Write([]byte(strconv.FormatFloat(m5.nice, 'g', 3, 64)))
	w.Write(crlf)
	w.Write([]byte("node_cpu_nice_ratio{duration=\"15s\"} "))
	w.Write([]byte(strconv.FormatFloat(m15.nice, 'g', 3, 64)))
	w.Write(crlf)

	w.Write([]byte("node_cpu_system_ratio{duration=\"1s\"} "))
	w.Write([]byte(strconv.FormatFloat(m1.system, 'g', 3, 64)))
	w.Write(crlf)
	w.Write([]byte("node_cpu_system_ratio{duration=\"5s\"} "))
	w.Write([]byte(strconv.FormatFloat(m5.system, 'g', 3, 64)))
	w.Write(crlf)
	w.Write([]byte("node_cpu_system_ratio{duration=\"15s\"} "))
	w.Write([]byte(strconv.FormatFloat(m15.system, 'g', 3, 64)))
	w.Write(crlf)

	w.Write([]byte("node_cpu_iowait_ratio{duration=\"1s\"} "))
	w.Write([]byte(strconv.FormatFloat(m1.iowait, 'g', 3, 64)))
	w.Write(crlf)
	w.Write([]byte("node_cpu_iowait_ratio{duration=\"5s\"} "))
	w.Write([]byte(strconv.FormatFloat(m5.iowait, 'g', 3, 64)))
	w.Write(crlf)
	w.Write([]byte("node_cpu_iowait_ratio{duration=\"15s\"} "))
	w.Write([]byte(strconv.FormatFloat(m15.iowait, 'g', 3, 64)))
	w.Write(crlf)

	w.Write([]byte("node_cpu_steal_ratio{duration=\"1s\"} "))
	w.Write([]byte(strconv.FormatFloat(m1.steal, 'g', 3, 64)))
	w.Write(crlf)
	w.Write([]byte("node_cpu_steal_ratio{duration=\"5s\"} "))
	w.Write([]byte(strconv.FormatFloat(m5.steal, 'g', 3, 64)))
	w.Write(crlf)
	w.Write([]byte("node_cpu_steal_ratio{duration=\"15s\"} "))
	w.Write([]byte(strconv.FormatFloat(m15.steal, 'g', 3, 64)))
	w.Write(crlf)

	w.Write([]byte("node_cpu_idle_ratio{duration=\"1s\"} "))
	w.Write([]byte(strconv.FormatFloat(m1.idle, 'g', 3, 64)))
	w.Write(crlf)
	w.Write([]byte("node_cpu_idle_ratio{duration=\"5s\"} "))
	w.Write([]byte(strconv.FormatFloat(m5.idle, 'g', 3, 64)))
	w.Write(crlf)
	w.Write([]byte("node_cpu_idle_ratio{duration=\"15s\"} "))
	w.Write([]byte(strconv.FormatFloat(m15.idle, 'g', 3, 64)))
	w.Write(crlf)

	for dev, io := range m1.devices {
		w.Write([]byte("node_disk_writes_bytes_per_second{duration=\"1s\", device=\"" + dev + "\"} "))
		w.Write([]byte(strconv.FormatFloat(io.bytes_write, 'g', 3, 64)))
		w.Write(crlf)
		w.Write([]byte("node_disk_reads_bytes_per_second{duration=\"1s\", device=\"" + dev + "\"} "))
		w.Write([]byte(strconv.FormatFloat(io.bytes_read, 'g', 3, 64)))
		w.Write(crlf)
	}

	for dev, io := range m5.devices {
		w.Write([]byte("node_disk_writes_bytes_per_second{duration=\"5s\", device=\"" + dev + "\"} "))
		w.Write([]byte(strconv.FormatFloat(io.bytes_write, 'g', 3, 64)))
		w.Write(crlf)
		w.Write([]byte("node_disk_reads_bytes_per_second{duration=\"5s\", device=\"" + dev + "\"} "))
		w.Write([]byte(strconv.FormatFloat(io.bytes_read, 'g', 3, 64)))
		w.Write(crlf)
	}

	for dev, io := range m15.devices {
		w.Write([]byte("node_disk_writes_bytes_per_second{duration=\"15s\", device=\"" + dev + "\"} "))
		w.Write([]byte(strconv.FormatFloat(io.bytes_write, 'g', 3, 64)))
		w.Write(crlf)
		w.Write([]byte("node_disk_reads_bytes_per_second{duration=\"15s\", device=\"" + dev + "\"} "))
		w.Write([]byte(strconv.FormatFloat(io.bytes_read, 'g', 3, 64)))
		w.Write(crlf)
	}
}

func getSystemMetrics(w http.ResponseWriter) {
	var crlf = []byte("\r\n")
	v, _ := mem.VirtualMemory()
	w.Write([]byte("node_memory_bytes_total "))
	w.Write([]byte(strconv.FormatUint(v.Total, 10)))
	w.Write(crlf)
	w.Write([]byte("node_memory_bytes_used "))
	w.Write([]byte(strconv.FormatUint(v.Used, 10)))
	w.Write(crlf)
	w.Write([]byte("node_memory_ratio_used "))
	w.Write([]byte(strconv.FormatFloat(v.UsedPercent/100.0, 'g', 3, 64)))
	w.Write(crlf)

	l, _ := load.Avg()
	w.Write([]byte("node_load{duration=\"1s\"} "))
	w.Write([]byte(strconv.FormatFloat(l.Load1, 'e', 3, 64)))
	w.Write(crlf)
	w.Write([]byte("node_load{duration=\"5s\"} "))
	w.Write([]byte(strconv.FormatFloat(l.Load5, 'e', 3, 64)))
	w.Write(crlf)
	w.Write([]byte("node_load{duration=\"15s\"} "))
	w.Write([]byte(strconv.FormatFloat(l.Load15, 'e', 3, 64)))
	w.Write(crlf)
}

func getConsulMetrics(w http.ResponseWriter) {
	metrics, err := loadConsulMetrics()
	if err != nil {
		panic(err)
	}
	c, ok := metrics.Counters["consul.client.rpc"]
	if ok {
		line := "consul_client_rpc_requests_total " + strconv.FormatUint(c.Count, 10) + "\r\n"
		w.Write([]byte(line))
	}
	c, ok = metrics.Counters["consul.rpc.request"]
	if ok {
		line := "consul_rpc_requests_total " + strconv.FormatUint(c.Count, 10) + "\r\n"
		w.Write([]byte(line))
	}
	c, ok = metrics.Counters["consul.raft.state.leader"]
	if ok {
		line := "consul_election_wins_total " + strconv.FormatUint(c.Count, 10) + "\r\n"
		w.Write([]byte(line))
	}
	c, ok = metrics.Counters["consul.raft.state.candidate"]
	if ok {
		line := "consul_election_total " + strconv.FormatUint(c.Count, 10) + "\r\n"
		w.Write([]byte(line))
	}
	c, ok = metrics.Counters["consul.raft.apply"]
	if ok {
		applypersec := float64(c.Count) / 10.0
		line := "consul_raft_apply_seconds " + strconv.FormatFloat(applypersec, 'g', 3, 64) + "\r\n"
		w.Write([]byte(line))
	}
	s, ok := metrics.Samples["consul.raft.commitTime"]
	if ok {
		line := "consul_raft_committime_seconds " + strconv.FormatFloat(s.Mean/1000, 'g', 3, 64) + "\r\n"
		w.Write([]byte(line))
	}
	s, ok = metrics.Samples["consul.raft.replication.appendEntries"]
	if ok {
		line := "consul_raft_replication_seconds " + strconv.FormatFloat(s.Mean/1000, 'g', 3, 64) + "\r\n"
		w.Write([]byte(line))
	}
}

func getNomadMetrics(w http.ResponseWriter) {
	url := "http://127.0.0.1:4646/v1/metrics?format=prometheus"
	response, err := http.Get(url)
	defer response.Body.Close()
	if err != nil {
		panic(err)
	}
	var crlf = []byte("\r\n")

	hostname, err := os.Hostname()
	if err != nil {
		panic(err)
	}
	hostname = strings.Replace(hostname, ".", "_", -1) + "_"
	scanner := bufio.NewScanner(response.Body)
	for scanner.Scan() {
		var line = scanner.Text()
		if strings.HasPrefix(line, "nomad") {
			var newLine = strings.Replace(line, hostname, "", 1)
			if strings.HasPrefix(newLine, "nomad_client") || strings.HasPrefix(newLine, "nomad_runtime") {
				w.Write([]byte(newLine))
				w.Write(crlf)
			}
		}
	}
	return
}
