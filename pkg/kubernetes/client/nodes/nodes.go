package nodes

import (
	cvm "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/cvm/v20170312"
)

var Nodes = make(map[string]string)

func AddNodeFromInstanceIdSet(res *cvm.RunInstancesResponse, clusterName string) {
	for _, id := range res.Response.InstanceIdSet {
		Nodes[*id] = clusterName
	}
}

func DeleteNodeFromInstanceIdSet(res *cvm.RunInstancesResponse) {
	for _, id := range res.Response.InstanceIdSet {
		delete(Nodes, *id)
	}
}

func GetClusterNameFromInstanceIdSet(res *cvm.RunInstancesResponse) string {
	return Nodes[*res.Response.InstanceIdSet[0]]
}
