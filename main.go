package main

import (
	"mobingi/ocean/app"
	"mobingi/ocean/pkg/kubernetes/client"
	"mobingi/ocean/pkg/log"
	"mobingi/ocean/pkg/storage"
)

func main() {
	// TODO Move to main init func
	storage.NewMongoClient()

	err := client.InitClustersAndNodes()
	if err != nil {
		log.Error(err)
		return
	}
	err = client.ClustersMonitor()
	if err != nil {
		log.Error(err)
		return
	}

	// tencent.Init()
	// client := &tencent.InstanceTencent{}

	// res, err := client.CreateInstance(1)
	// if err != nil {
	// 	panic(err)
	// }
	// fmt.Println(res.ToJsonString())
	// return

	// filters := []*batch.Filter{&batch.Filter{
	// 	Name:   common.StringPtr("env-name"),
	// 	Values: []*string{common.StringPtr(tencent.EnvName)},
	// }}
	// response, err := client.DescribeComputeEnvs(filters)
	// if err != nil {
	// 	panic(err)
	// }
	// if *response.Response.TotalCount > 0 {
	// 	fmt.Printf("已存在计算环境:%s \n", tencent.EnvName)
	// } else {
	// 	env, err := client.CreateComputeEnv()
	// 	if err != nil {
	// 		fmt.Println(err.Error())
	// 	}
	// 	fmt.Println(env.ToJsonString())
	// }

	if err := app.Start(); err != nil {
		panic(err)
	}
}
