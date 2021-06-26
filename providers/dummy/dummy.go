package dummy

import (
	"dimitrios_liappis/instancepeek/model"
	"fmt"
	"math/rand"
	"strings"
	"time"
)

func Convert(instanceTypes []string) ([]model.InstanceInfo, error) {
	rand.Seed(time.Now().UnixNano())

	var instanceInfos []model.InstanceInfo
	for _, instanceType := range instanceTypes {
		instanceInfos = append(instanceInfos, model.InstanceInfo{
			InstanceType: model.Data{
				Label: "Instance Type",
				Value: instanceType,
			},
			CPU: model.CPUInfo{
				VCPUCount: model.Data{
					Label: "vCPUs",
					Value: fmt.Sprint(rand.Intn(36)),
				},
			},
			Disk: model.LocalDiskInfo{
				Typ: model.Data{
					Label: "Local Disk Type",
					Value: randString(),
				},
				Count: model.Data{
					Label: "Local Disk Count",
					Value: fmt.Sprint(rand.Intn(24)),
				},
				SizeGiB: model.Data{
					Label: "Local Disk Size(GiB)",
					Value: fmt.Sprint(rand.Intn(65535)),
				},
			},
			Memory: model.MemoryInfo{
				SizeMiB: model.Data{
					Label: "Memory MiB",
					Value: fmt.Sprint(rand.Intn(262144)),
				},
			},
			Network: model.NetworkInfo{
				Performance: model.Data{
					Label: "Network Speed",
					Value: "up to " + fmt.Sprint(rand.Intn(100)) + "Gbps",
				},
			},
			Hardware: model.VMInfo{
				Hypervisor: model.Data{
					Label: "Hypervisor",
					Value: randString(),
				},
				Baremetal: model.Data{
					Label: "Baremetal",
					Value: fmt.Sprint(rand.Intn(2)),
				},
			},
			Meta: map[string]string{
				"EBSBaselineIops":             "18750",
				"EBSBaselineThroughputInMBps": "593.75",
				"EBSMaximumIops":              "18750",
				"EBSMaximumThroughputInMBps":  "593.75",
			},
		})

		time.Sleep(time.Millisecond * time.Duration(rand.Float32()*500))
	}
	return instanceInfos, nil
}

func randString() string {
	var alphabet = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ")
	var str strings.Builder
	for i := 0; i < rand.Intn(len(alphabet)); i++ {
		randRune := alphabet[rand.Intn(len(alphabet))]
		str.WriteRune(randRune)
	}
	return str.String()
}
