package main

import (
	"fmt"
	"mobingi/ocean/app"
	"mobingi/ocean/pkg/services/tencent"
	"mobingi/ocean/pkg/storage"

	batch "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/batch/v20170312"
	"github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/common"
)

func main() {
	// TODO Move to main init func
	storage.NewMongoClient()

	tencent.Init()
	client := &tencent.InstanceTencent{}

	filters := []*batch.Filter{&batch.Filter{
		Name:   common.StringPtr("env-name"),
		Values: []*string{common.StringPtr(tencent.EnvName)},
	}}
	response, err := client.DescribeComputeEnvs(filters)
	if err != nil {
		panic(err)
	}
	if *response.Response.TotalCount > 0 {
		fmt.Printf("已存在计算环境:%s \n", tencent.EnvName)
	} else {
		env, err := client.CreateComputeEnv()
		if err != nil {
			fmt.Println(err.Error())
		}
		fmt.Println(env.ToJsonString())
	}

	// storage := storage.NewStorage()
	// kubeconfig, err := storage.GetKubeconf("kubernetes", "admin.conf")
	// if err != nil {
	// 	log.Error(err)
	// }
	// err = client.Init(kubeconfig)
	// if err != nil {
	// 	log.Error(err)
	// }

	// nodes, err := client.GetNode()
	// fmt.Println(nodes.Items)
	if err := app.Start(); err != nil {
		panic(err)
	}
}
