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

// Converts human friendly names of instance types to native instance types as expected by the AWS SDK
func toNativeInstanceTypes(instanceTypes []string) []types.InstanceType {
	var nativeInstanceTypes []types.InstanceType
	for _, it := range instanceTypes {
		nativeInstanceTypes = append(nativeInstanceTypes, types.InstanceType(it))
	}
	return nativeInstanceTypes
}

// Converts AWS InstanceTypes to a standardized model
func Convert(ctx context.Context, instanceTypes []string) ([]model.InstanceInfo, error) {
	return ConfigurableConvert(ctx, instanceTypes, Client("us-east-2"))
}

// Testable version of Convert allows client mocking
func ConfigurableConvert(ctx context.Context, instanceTypes []string, api EC2DescribeInstanceTypesAPI) ([]model.InstanceInfo, error) {
	input := &ec2.DescribeInstanceTypesInput{
		InstanceTypes: toNativeInstanceTypes(instanceTypes),
	}

	resp, err := api.DescribeInstanceTypes(ctx, input)
	if err != nil {
		panic("Error " + err.Error())
	}

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
