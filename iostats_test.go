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
	"testing"
	log "github.com/sirupsen/logrus"
	"time"
)

func TestIoStats(t *testing.T) {
	log.SetLevel(log.InfoLevel)
	println("Running stats in goroutine")
	go RunStats()
	timeChan := time.Tick(5 * time.Second)
	for range timeChan {
		m1, m5, m15 := GetStats()
		log.Printf("1s Average: %v \n", m1)
		log.Printf("5s Average: %v \n", m5)
		log.Printf("15s Average: %v \n", m15)
	}
}
