package tencent

import (
	"fmt"
	"os"
	"time"

	cvm "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/cvm/v20170312"

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
	ImageName    = "img-lggrctuj"
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
			ImageId:      common.StringPtr(ImageName),
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

func (c *InstanceTencent) CreateSpotInstance(number int64) (*cvm.RunInstancesResponse, error) {
	cpf := profile.NewClientProfile()
	client, _ := cvm.NewClient(credential, regions.Beijing, cpf)
	request := cvm.NewRunInstancesRequest()
	request.Placement = &cvm.Placement{
		Zone: common.StringPtr("ap-beijing-1"),
	}
	request.ImageId = common.StringPtr(ImageName)
	request.InstanceChargeType = common.StringPtr("SPOTPAID")
	request.InstanceType = common.StringPtr("S2.SMALL1")
	request.SystemDisk = &cvm.SystemDisk{
		DiskType: common.StringPtr("CLOUD_PREMIUM"),
		DiskSize: common.Int64Ptr(50),
	}
	request.InternetAccessible = &cvm.InternetAccessible{
		InternetChargeType:      common.StringPtr("TRAFFIC_POSTPAID_BY_HOUR"),
		InternetMaxBandwidthOut: common.Int64Ptr(1),
		PublicIpAssigned:        common.BoolPtr(true),
	}
	request.InstanceCount = common.Int64Ptr(number)
	request.InstanceName = common.StringPtr(InstanceName)
	request.LoginSettings = &cvm.LoginSettings{
		Password: common.StringPtr("A!Y947337"),
	}
	request.ClientToken = common.StringPtr(time.Now().String())
	request.HostName = common.StringPtr("k8s001")
	request.InstanceMarketOptions = &cvm.InstanceMarketOptionsRequest{
		MarketType: common.StringPtr("spot"),
		SpotOptions: &cvm.SpotMarketOptions{
			MaxPrice:         common.StringPtr("0.03"),
			SpotInstanceType: common.StringPtr("one-time"),
		},
	}

	res, err := client.RunInstances(request)
	if err != nil {
		return nil, err
	}
	return res, nil
}

func (c *InstanceTencent) CreateCommonInstance(number int64) (*cvm.RunInstancesResponse, error) {
	cpf := profile.NewClientProfile()
	client, _ := cvm.NewClient(credential, regions.Beijing, cpf)
	request := cvm.NewRunInstancesRequest()
	request.Placement = &cvm.Placement{
		Zone: common.StringPtr("ap-beijing-1"),
	}
	request.ImageId = common.StringPtr(ImageName)
	request.InstanceChargeType = common.StringPtr("POSTPAID_BY_HOUR")
	request.InstanceType = common.StringPtr("S2.SMALL1")
	request.SystemDisk = &cvm.SystemDisk{
		DiskType: common.StringPtr("CLOUD_PREMIUM"),
		DiskSize: common.Int64Ptr(50),
	}
	request.InternetAccessible = &cvm.InternetAccessible{
		InternetChargeType:      common.StringPtr("TRAFFIC_POSTPAID_BY_HOUR"),
		InternetMaxBandwidthOut: common.Int64Ptr(1),
		PublicIpAssigned:        common.BoolPtr(true),
	}
	request.InstanceCount = common.Int64Ptr(number)
	request.InstanceName = common.StringPtr(InstanceName)
	request.LoginSettings = &cvm.LoginSettings{
		Password: common.StringPtr("A!Y947337"),
	}
	request.ClientToken = common.StringPtr(time.Now().String())

	res, err := client.RunInstances(request)
	if err != nil {
		return nil, err
	}
	return res, nil
}

func (c *InstanceTencent) CreateInstance(number int64) (*cvm.RunInstancesResponse, error) {
	// 是否分批次创建更好
	res, err := c.CreateSpotInstance(number)
	if err != nil {
		// TODO 普通实例应自动替换为spot实例
		res, err = c.CreateCommonInstance(number)
		if err != nil {
			return nil, err
		}
	}
	return res, nil
}
