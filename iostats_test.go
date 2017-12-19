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
