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
	getNomadMetrics(w)
}

func getIoStats(w http.ResponseWriter) {
	m1, m5, m15 := GetStats()
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
	// TODO: add device IO
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
