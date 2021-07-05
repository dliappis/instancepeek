package aws

import (
	"context"
	"dimitrios_liappis/instancepeek/model"
	"encoding/json"
	"io/ioutil"
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

func TestConvert(t *testing.T) {
	b, err := ioutil.ReadFile("instancetypeinfo.json") // file created using respDecrypted, _ := json.MarshalIndent(resp, "", "\t"); fmt.Println(string(respDecrypted)) in aws.ConfigurableConvert()
	if err != nil {
		t.Fatal(err)
	}
	var testInstanceType types.InstanceTypeInfo
	json.Unmarshal(b, &testInstanceType)

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
						InstanceTypes: []types.InstanceTypeInfo{
							testInstanceType,
						},
					}, nil
				})
			},
			expect: []model.InstanceInfo{
				{
					InstanceType: model.Data{
						Label: "Instance Type",
						Value: "m5d.4xlarge",
					},
					CPU: model.CPUInfo{
						VCPUCount: model.Data{
							Label: "vCPUs",
							Value: "16",
						},
					},
					Disk: model.LocalDiskInfo{
						Typ: model.Data{
							Label: "Local Disk Type",
							Value: "ssd",
						},
						Count: model.Data{
							Label: "Local Disk Count",
							Value: "2",
						},
						SizeGiB: model.Data{
							Label: "Local Disk Size(GB)",
							Value: "300",
						},
					},
					Memory: model.MemoryInfo{
						SizeMiB: model.Data{
							Label: "Memory MiB",
							Value: "65536",
						},
					},
					Network: model.NetworkInfo{
						Performance: model.Data{
							Label: "Network Speed",
							Value: "Up to 10 Gigabit",
						},
					},
					Hardware: model.VMInfo{
						Hypervisor: model.Data{
							Label: "Hypervisor",
							Value: "nitro",
						},
						Baremetal: model.Data{
							Label: "Baremetal",
							Value: "false",
						},
					},
					Meta: map[string]string{
						"EBSBaselineIops":              "18750",
						"EBSBaselineThroughput (MBps)": "593.75",
						"EBSMaxIops":                   "18750",
						"EBSMaxThroughput (MBps)":      "593.75",
					},
				},
			},
		},
	}

	for i, tt := range cases {
		t.Run(strconv.Itoa(i), func(t *testing.T) {
			ctx := context.TODO()
			content, err := ConfigurableConvert(ctx, []string{"m5d.4xlarge"}, tt.client(t))
			if err != nil {
				t.Fatalf("expect no error, got %v", err)
			}
			if e, a := tt.expect, content; !cmp.Equal(e, a) {
				t.Errorf("expect %v, got %v", e, a)
			}
		})
	}
}
