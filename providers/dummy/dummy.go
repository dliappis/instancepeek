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
			InstanceType: instanceType,
			CPU: model.CPUInfo{
				VCPUCount: fmt.Sprint(rand.Intn(36)),
			},
			Disk: model.LocalDiskInfo{
				Typ:     randString(),
				Count:   fmt.Sprint(rand.Intn(24)),
				SizeGiB: fmt.Sprint(rand.Intn(65535)),
			},
			Memory: model.MemoryInfo{
				SizeMiB: fmt.Sprint(rand.Intn(262144)),
			},
			Network: model.NetworkInfo{
				Performance: "up to " + fmt.Sprint(rand.Intn(100)) + "Gbps",
			},
			Hardware: model.VMInfo{
				Hypervisor: randString(),
				Baremetal:  fmt.Sprint(rand.Intn(2)),
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
