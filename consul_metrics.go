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
	"encoding/json"
	"io/ioutil"
)

type metricsJson struct {
	Gauges   []Gauge
	Counters []Counter
	Samples  []Sample
}

type ConsulMetrics struct {
	Gauges   map[string]Gauge
	Counters map[string]Counter
	Samples  map[string]Sample
}

type Sample struct {
	Name   string
	Count  uint64
	Sum    float64
	Min    float64
	Max    float64
	Mean   float64
	Stddev float64
}

type Gauge struct {
	Name  string
	Value uint64
}

type Counter struct {
	Name   string
	Count  uint64
	Sum    uint64
	Min    uint64
	Max    uint64
	Mean   float64
	Stddev float64
}

func loadConsulMetrics() (*ConsulMetrics, error) {
	url := "http://127.0.0.1:8500/v1/agent/metrics"
	response, err := http.Get(url)
	defer response.Body.Close()
	if err != nil {
		return nil, err
	}
	met := metricsJson{}
	body, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return nil, err
	}
	err = json.Unmarshal(body, &met)
	res := ConsulMetrics{
		Gauges:   make(map[string]Gauge),
		Counters: make(map[string]Counter),
		Samples:  make(map[string]Sample),
	}
	for _,g:= range met.Gauges {
		res.Gauges[g.Name]=g
	}
	for _,c:= range met.Counters{
		res.Counters[c.Name]=c
	}
	for _,s:= range met.Samples{
		res.Samples[s.Name]=s
	}
	return &res, err
}
