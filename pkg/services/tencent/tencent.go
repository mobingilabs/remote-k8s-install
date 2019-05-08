package tencent

import (
	"fmt"
	"os"

	batch "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/batch/v20170312"
	"github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/common"
	"github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/common/errors"
	"github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/common/profile"
	"github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/common/regions"
)

type InstanceTencent struct{}

const (
	EnvName      = "kubernetes"
	InstanceName = "kubernetes-node"
	ClientToken  = "create-compute-env"
)

var credential *common.Credential

func Init() {
	credential = common.NewCredential(os.Getenv("secretID"), os.Getenv("secretKey"))
}

func (c *InstanceTencent) CreateComputeEnv() (*batch.CreateComputeEnvResponse, error) {
	cpf := profile.NewClientProfile()
	client, _ := batch.NewClient(credential, regions.Beijing, cpf)
	request := batch.NewCreateComputeEnvRequest()
	request.Placement = &batch.Placement{
		Zone: common.StringPtr("ap-beijing-1"),
	}
	request.ClientToken = common.StringPtr(ClientToken)
	request.ComputeEnv = &batch.NamedComputeEnv{
		EnvName:                 common.StringPtr(EnvName),
		DesiredComputeNodeCount: common.Int64Ptr(1),
		EnvDescription:          common.StringPtr("env description"),
		EnvType:                 common.StringPtr("MANAGED"),
		EnvData: &batch.EnvData{
			InstanceType: common.StringPtr("S2.SMALL1"),
			ImageId:      common.StringPtr("img-8uwydzx3"),
			SystemDisk: &batch.SystemDisk{
				DiskType: common.StringPtr("CLOUD_PREMIUM"),
				DiskSize: common.Int64Ptr(50),
			},
			InternetAccessible: &batch.InternetAccessible{
				InternetChargeType:      common.StringPtr("TRAFFIC_POSTPAID_BY_HOUR"),
				InternetMaxBandwidthOut: common.Int64Ptr(1),
				PublicIpAssigned:        common.BoolPtr(true),
			},
			InstanceName: common.StringPtr(InstanceName),
			LoginSettings: &batch.LoginSettings{
				Password: common.StringPtr("A!Y947337"),
			},
			InstanceChargeType: common.StringPtr("SPOTPAID"),
			InstanceMarketOptions: &batch.InstanceMarketOptionsRequest{
				MarketType: common.StringPtr("spot"),
				SpotOptions: &batch.SpotMarketOptions{
					MaxPrice:         common.StringPtr("0.03"),
					SpotInstanceType: common.StringPtr("one-time"),
				},
			},
		},
	}

	response, err := client.CreateComputeEnv(request)
	if _, ok := err.(*errors.TencentCloudSDKError); ok {
		return nil, fmt.Errorf("An API error has returned: %s", err)
	}
	if err != nil {
		return nil, err
	}
	return response, nil
}

func (c *InstanceTencent) DescribeComputeEnvs(filters []*batch.Filter) (*batch.DescribeComputeEnvsResponse, error) {
	cpf := profile.NewClientProfile()
	client, _ := batch.NewClient(credential, regions.Beijing, cpf)
	request := batch.NewDescribeComputeEnvsRequest()
	request.Filters = filters

	response, err := client.DescribeComputeEnvs(request)
	if _, ok := err.(*errors.TencentCloudSDKError); ok {
		return nil, fmt.Errorf("An API error has returned: %s", err)
	}
	if err != nil {
		return nil, err
	}
	return response, nil
}

func (c *InstanceTencent) DescribeComputeEnv(envId *string) (*batch.DescribeComputeEnvResponse, error) {
	cpf := profile.NewClientProfile()
	client, _ := batch.NewClient(credential, regions.Beijing, cpf)
	request := batch.NewDescribeComputeEnvRequest()
	request.EnvId = envId

	response, err := client.DescribeComputeEnv(request)
	if _, ok := err.(*errors.TencentCloudSDKError); ok {
		return nil, fmt.Errorf("An API error has returned: %s", err)
	}
	if err != nil {
		return nil, err
	}
	return response, nil
}
