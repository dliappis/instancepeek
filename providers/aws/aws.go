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
		baseInstanceInfo := model.InstanceInfo{
			InstanceType: model.Data{
				Label: "Instance Type",
				Value: string(it.InstanceType),
			},
			CPU: model.CPUInfo{
				VCPUCount: model.Data{
					Label: "vCPUs",
					Value: fmt.Sprint(*it.VCpuInfo.DefaultVCpus),
				},
			},
			Memory: model.MemoryInfo{
				SizeMiB: model.Data{
					Label: "Memory MiB",
					Value: fmt.Sprint(*it.MemoryInfo.SizeInMiB),
				},
			},
			Network: model.NetworkInfo{
				Performance: model.Data{
					Label: "Network Speed",
					Value: string(*it.NetworkInfo.NetworkPerformance),
				},
			},
			Hardware: model.VMInfo{
				Hypervisor: model.Data{
					Label: "Hypervisor",
					Value: string(it.Hypervisor),
				},
				Baremetal: model.Data{
					Label: "Baremetal",
					Value: strconv.FormatBool(*it.BareMetal),
				},
			},
			Meta: map[string]string{
				"EBSBaselineIops":              fmt.Sprint(*it.EbsInfo.EbsOptimizedInfo.BaselineIops),
				"EBSBaselineThroughput (MBps)": fmt.Sprint(*it.EbsInfo.EbsOptimizedInfo.BaselineThroughputInMBps),
				"EBSMaxIops":                   fmt.Sprint(*it.EbsInfo.EbsOptimizedInfo.MaximumIops),
				"EBSMaxThroughput (MBps)":      fmt.Sprint(*it.EbsInfo.EbsOptimizedInfo.MaximumThroughputInMBps),
			},
		}

		if it.InstanceStorageInfo != nil {
			localDiskInfo := model.LocalDiskInfo{
				Typ: model.Data{
					Label: "Local Disk Type",
					Value: string(it.InstanceStorageInfo.Disks[0].Type),
				},
				Count: model.Data{
					Label: "Local Disk Count",
					Value: fmt.Sprint(*it.InstanceStorageInfo.Disks[0].Count),
				},
				SizeGiB: model.Data{
					Label: "Local Disk Size(GB)",
					Value: fmt.Sprint(*it.InstanceStorageInfo.Disks[0].SizeInGB),
				},
			}
			baseInstanceInfo.Disk = localDiskInfo
		}

		instanceInfos = append(instanceInfos, baseInstanceInfo)
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
