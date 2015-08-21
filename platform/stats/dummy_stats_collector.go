package stats

import (
	"time"
)

type dummyStatsCollector struct{}

func NewDummyStatsCollector() (collector Collector) {
	return dummyStatsCollector{}
}

func (p dummyStatsCollector) StartCollecting(collectionInterval time.Duration, latestGotUpdated chan struct{}) {
}

func (p dummyStatsCollector) GetCPULoad() (load CPULoad, err error) {
	return
}

func (p dummyStatsCollector) GetCPUStats() (stats CPUStats, err error) {
	stats.Total = 1
	return
}

func (p dummyStatsCollector) GetMemStats() (usage Usage, err error) {
	usage.Total = 1
	return
}

func (p dummyStatsCollector) GetSwapStats() (usage Usage, err error) {
	usage.Total = 1
	return
}

func (p dummyStatsCollector) GetDiskStats(devicePath string) (stats DiskStats, err error) {
	stats.DiskUsage.Total = 1
	stats.InodeUsage.Total = 1
	return
}

func (p dummyStatsCollector) GetProcessStats() (stats []ProcessStat, err error) {
	stats = []ProcessStat{
		ProcessStat{
			Name:  "cloud_controller",
			State: "running",
		},
		ProcessStat{
			Name:  "cloud_controller_worker",
			State: "running",
		},
	}
	return
}
