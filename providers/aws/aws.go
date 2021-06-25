// Builds a catalog of AWS instance type details
package aws

import (
	"context"
	"dimitrios_liappis/instancepeek/model"
	"fmt"
	"strconv"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/ec2"
	"github.com/aws/aws-sdk-go-v2/service/ec2/types"
)

type EC2DescribeInstanceTypesAPI interface {
	DescribeInstanceTypes(ctx context.Context, params *ec2.DescribeInstanceTypesInput, optFns ...func(*ec2.Options)) (*ec2.DescribeInstanceTypesOutput, error)
}

// Converts AWS InstanceTypes to my own type
func Convert(ctx context.Context, instanceTypes []string, api EC2DescribeInstanceTypesAPI) ([]model.InstanceInfo, error) {
	var iTs []types.InstanceType
	for _, iT := range instanceTypes {
		iTs = append(iTs, types.InstanceType(iT))
	}
	input := &ec2.DescribeInstanceTypesInput{
		InstanceTypes: iTs,
	}

	resp, err := api.DescribeInstanceTypes(ctx, input)
	if err != nil {
		panic("Error " + err.Error())
	}

	// respDecrypted, _ := json.MarshalIndent(resp, "", "\t")
	// fmt.Println(string(respDecrypted))
	var instanceInfos []model.InstanceInfo
	for _, it := range resp.InstanceTypes {
		instanceInfos = append(instanceInfos, model.InstanceInfo{
			InstanceType: string(it.InstanceType),
			CPU: model.CPUInfo{
				VCPUCount: fmt.Sprint(*it.VCpuInfo.DefaultVCpus),
			},
			Disk: model.LocalDiskInfo{
				Typ:     string(it.InstanceStorageInfo.Disks[0].Type),
				Count:   fmt.Sprint(*it.InstanceStorageInfo.Disks[0].Count),
				SizeGiB: fmt.Sprint(*it.InstanceStorageInfo.Disks[0].SizeInGB),
			},
			Memory: model.MemoryInfo{
				SizeMiB: fmt.Sprint(*it.MemoryInfo.SizeInMiB),
			},
			Network: model.NetworkInfo{
				Performance: string(*it.NetworkInfo.NetworkPerformance),
			},
			Hardware: model.VMInfo{
				Hypervisor: string(it.Hypervisor),
				Baremetal:  strconv.FormatBool(*it.BareMetal),
			},
		})

	}
	return instanceInfos, nil
}

func baseCfg(region string) aws.Config {
	// gets the AWS credentials from the default file or from the EC2 instance profile
	cfg, err := config.LoadDefaultConfig(context.TODO(),
		config.WithRegion(region))
	if err != nil {
		panic("Unable to load AWS SDK config: " + err.Error())
	}

	return cfg
}

// Sets up the client for API operations
func Client(region string) *ec2.Client {
	cfg := baseCfg(region)

	client := ec2.NewFromConfig(cfg)
	return client
}
