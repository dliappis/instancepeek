package aws

import (
	"context"
	"dimitrios_liappis/instancepeek/model"
	"strconv"
	"testing"

	"github.com/aws/aws-sdk-go-v2/service/ec2"
	"github.com/aws/aws-sdk-go-v2/service/ec2/types"
	"github.com/google/go-cmp/cmp"
)

type mockDescriptInstanceTypesAPI func(ctx context.Context, params *ec2.DescribeInstanceTypesInput, optFns ...func(*ec2.Options)) (*ec2.DescribeInstanceTypesOutput, error)

func (m mockDescriptInstanceTypesAPI) DescribeInstanceTypes(ctx context.Context, params *ec2.DescribeInstanceTypesInput, optFns ...func(*ec2.Options)) (*ec2.DescribeInstanceTypesOutput, error) {
	return m(ctx, params, optFns...)
}

func newBoolPtr(val bool) *bool {
	return &val
}

func newStringPtr(val string) *string {
	return &val
}

func newInt32Ptr(val int32) *int32 {
	return &val
}

func newInt64Ptr(val int64) *int64 {
	return &val
}

func newFloat64Ptr(val float64) *float64 {
	return &val
}

var InstanceTypesResponse = []types.InstanceTypeInfo{
	{
		AutoRecoverySupported:         newBoolPtr(false),
		BareMetal:                     newBoolPtr(false),
		BurstablePerformanceSupported: newBoolPtr(false),
		CurrentGeneration:             newBoolPtr(true),
		DedicatedHostsSupported:       newBoolPtr(true),
		EbsInfo: &types.EbsInfo{
			EbsOptimizedInfo: &types.EbsOptimizedInfo{
				BaselineBandwidthInMbps:  newInt32Ptr(4750),
				BaselineIops:             newInt32Ptr(18750),
				BaselineThroughputInMBps: newFloat64Ptr(593.75),
				MaximumBandwidthInMbps:   newInt32Ptr(4750),
				MaximumIops:              newInt32Ptr(18750),
				MaximumThroughputInMBps:  newFloat64Ptr(593.75),
			},
			EbsOptimizedSupport: "default",
			EncryptionSupport:   "supported",
			NvmeSupport:         "required",
		},
		FpgaInfo:                 nil,
		FreeTierEligible:         newBoolPtr(false),
		GpuInfo:                  nil,
		HibernationSupported:     newBoolPtr(false),
		Hypervisor:               "nitro",
		InferenceAcceleratorInfo: nil,
		InstanceStorageInfo: &types.InstanceStorageInfo{
			Disks: []types.DiskInfo{
				{
					Count:    newInt32Ptr(2),
					SizeInGB: newInt64Ptr(300),
					Type:     "ssd",
				}},
			NvmeSupport:   "required",
			TotalSizeInGB: newInt64Ptr(600),
		},
		InstanceStorageSupported: newBoolPtr(true),
		InstanceType:             "m5d.4xlarge",
		MemoryInfo: &types.MemoryInfo{
			SizeInMiB: newInt64Ptr(65536),
		},
		NetworkInfo: &types.NetworkInfo{
			DefaultNetworkCardIndex:   newInt32Ptr(0),
			EfaInfo:                   nil,
			EfaSupported:              newBoolPtr(false),
			EnaSupport:                "required",
			Ipv4AddressesPerInterface: newInt32Ptr(30),
			Ipv6AddressesPerInterface: newInt32Ptr(30),
			Ipv6Supported:             newBoolPtr(true),
			MaximumNetworkCards:       newInt32Ptr(1),
			MaximumNetworkInterfaces:  newInt32Ptr(8),
			NetworkCards: []types.NetworkCardInfo{
				{
					MaximumNetworkInterfaces: newInt32Ptr(8),
					NetworkCardIndex:         newInt32Ptr(0),
					NetworkPerformance:       newStringPtr("Up to 10 Gigabit"),
				},
			},
			NetworkPerformance: newStringPtr("Up to 10 Gigabit"),
		},
		PlacementGroupInfo: &types.PlacementGroupInfo{
			SupportedStrategies: []types.PlacementGroupStrategy{"cluster", "partition", "spread"},
		},
		ProcessorInfo: &types.ProcessorInfo{
			SupportedArchitectures:   []types.ArchitectureType{"x86_64"},
			SustainedClockSpeedInGhz: newFloat64Ptr(3.1),
		},
		SupportedBootModes:           []types.BootModeType{"legacy-bios", "uefi"},
		SupportedRootDeviceTypes:     []types.RootDeviceType{"ebs"},
		SupportedUsageClasses:        []types.UsageClassType{"on-demand", "spot"},
		SupportedVirtualizationTypes: []types.VirtualizationType{"hvm"},
		VCpuInfo: &types.VCpuInfo{
			DefaultCores:          newInt32Ptr(8),
			DefaultThreadsPerCore: newInt32Ptr(2),
			DefaultVCpus:          newInt32Ptr(16),
			ValidCores:            []int32{2, 4, 6, 8},
			ValidThreadsPerCore:   []int32{1, 2},
		},
	},
}

func TestConvert(t *testing.T) {
	cases := []struct {
		client func(t *testing.T) EC2DescribeInstanceTypesAPI
		input  *ec2.DescribeInstanceTypesInput
		expect []model.InstanceInfo
	}{
		{
			client: func(t *testing.T) EC2DescribeInstanceTypesAPI {
				return mockDescriptInstanceTypesAPI(func(ctx context.Context, params *ec2.DescribeInstanceTypesInput, optFns ...func(*ec2.Options)) (*ec2.DescribeInstanceTypesOutput, error) {
					t.Helper()

					return &ec2.DescribeInstanceTypesOutput{
						InstanceTypes: InstanceTypesResponse,
					}, nil
				})
			},
			expect: []model.InstanceInfo{
				{
					InstanceType: "m5d.4xlarge",
					CPU: model.CPUInfo{
						VCPUCount: "16",
					},
					Disk: model.LocalDiskInfo{
						Typ:     "ssd",
						Count:   "2",
						SizeGiB: "300",
					},
					Memory: model.MemoryInfo{
						SizeMiB: "65536",
					},
					Network: model.NetworkInfo{
						Performance: "Up to 10 Gigabit",
					},
					Hardware: model.VMInfo{
						Hypervisor: "nitro",
						Baremetal:  "false",
					},
				},
			},
		},
	}

	for i, tt := range cases {
		t.Run(strconv.Itoa(i), func(t *testing.T) {
			ctx := context.TODO()
			content, err := Convert(ctx, []string{"m5d.4xlarge"}, tt.client(t))
			if err != nil {
				t.Fatalf("expect no error, got %v", err)
			}
			if e, a := tt.expect, content; !cmp.Equal(e, a) {
				t.Errorf("expect %v, got %v", e, a)
			}
		})
	}
}
