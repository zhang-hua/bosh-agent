package vitals_test

import (
	"time"
	"errors"

	. "github.com/cloudfoundry/bosh-agent/internal/github.com/onsi/ginkgo"
	. "github.com/cloudfoundry/bosh-agent/internal/github.com/onsi/gomega"

	boshassert "github.com/cloudfoundry/bosh-agent/internal/github.com/cloudfoundry/bosh-utils/assert"
	boshstats "github.com/cloudfoundry/bosh-agent/platform/stats"
	fakestats "github.com/cloudfoundry/bosh-agent/platform/stats/fakes"
	. "github.com/cloudfoundry/bosh-agent/platform/vitals"
	boshdirs "github.com/cloudfoundry/bosh-agent/settings/directories"
)

func buildVitalsService() (statsCollector *fakestats.FakeCollector, service Service) {
	dirProvider := boshdirs.NewProvider("/fake/base/dir")
	statsCollector = &fakestats.FakeCollector{
		CPULoad: boshstats.CPULoad{
			One:     0.2,
			Five:    4.55,
			Fifteen: 1.123,
		},
		StartCollectingCPUStats: boshstats.CPUStats{
			User:  56,
			Sys:   10,
			Wait:  1,
			Total: 100,
		},
		MemStats: boshstats.Usage{
			Used:  700 * 1024,
			Total: 1000 * 1024,
		},
		SwapStats: boshstats.Usage{
			Used:  600 * 1024,
			Total: 1000 * 1024,
		},
		DiskStats: map[string]boshstats.DiskStats{
			"/": boshstats.DiskStats{
				DiskUsage:  boshstats.Usage{Used: 100, Total: 200},
				InodeUsage: boshstats.Usage{Used: 50, Total: 500},
			},
			dirProvider.DataDir(): boshstats.DiskStats{
				DiskUsage:  boshstats.Usage{Used: 15, Total: 20},
				InodeUsage: boshstats.Usage{Used: 10, Total: 50},
			},
			dirProvider.StoreDir(): boshstats.DiskStats{
				DiskUsage:  boshstats.Usage{Used: 2, Total: 2},
				InodeUsage: boshstats.Usage{Used: 3, Total: 4},
			},
		},
		ProcessStats: []boshstats.ProcessStat{
			boshstats.ProcessStat{
				Name: "fake-process-1",
				State: "running",
			},
			boshstats.ProcessStat{
				Name: "fake-process-2",
				State: "initializing",
			},
		},
	}

	service = NewService(statsCollector, dirProvider)
	statsCollector.StartCollecting(1*time.Millisecond, nil)
	return
}
func init() {
	Describe("Testing with Ginkgo", func() {
		It("vitals construction", func() {
			_, service := buildVitalsService()
			vitals, err := service.Get()

			expectedVitals := map[string]interface{}{
				"cpu": map[string]string{
					"sys":  "10.0",
					"user": "56.0",
					"wait": "1.0",
				},
				"disk": map[string]interface{}{
					"system": map[string]string{
						"percent":       "50",
						"inode_percent": "10",
					},
					"ephemeral": map[string]string{
						"percent":       "75",
						"inode_percent": "20",
					},
					"persistent": map[string]string{
						"percent":       "100",
						"inode_percent": "75",
					},
				},
				"load": []string{"0.20", "4.55", "1.12"},
				"mem": map[string]string{
					"kb":      "700",
					"percent": "70",
				},
				"swap": map[string]string{
					"kb":      "600",
					"percent": "60",
				},
				"process": []interface{}{
					map[string]string{
						"name": "fake-process-1",
						"state": "running",
					},
					map[string]string{
						"name": "fake-process-2",
						"state": "initializing",
					},
				},
			}

			Expect(err).ToNot(HaveOccurred())

			boshassert.MatchesJSONMap(GinkgoT(), vitals, expectedVitals)
		})

		It("getting vitals when missing disks", func() {

			statsCollector, service := buildVitalsService()
			statsCollector.DiskStats = map[string]boshstats.DiskStats{
				"/": boshstats.DiskStats{
					DiskUsage:  boshstats.Usage{Used: 100, Total: 200},
					InodeUsage: boshstats.Usage{Used: 50, Total: 500},
				},
			}

			vitals, err := service.Get()
			Expect(err).ToNot(HaveOccurred())

			boshassert.LacksJSONKey(GinkgoT(), vitals.Disk, "ephemeral")
			boshassert.LacksJSONKey(GinkgoT(), vitals.Disk, "persistent")
		})

		It("get getting vitals on system disk error", func() {

			statsCollector, service := buildVitalsService()
			statsCollector.DiskStats = map[string]boshstats.DiskStats{}

			_, err := service.Get()
			Expect(err).To(HaveOccurred())
		})

		It("raise error when failing to get process stats", func() {
			statsCollector, service := buildVitalsService()
			statsCollector.ProcessStatsErr = errors.New("error on getting processes state")

			_, err := service.Get()
			Expect(err).To(HaveOccurred())
		})
	})
}
