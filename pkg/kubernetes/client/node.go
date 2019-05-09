package client

import (
	"mobingi/ocean/pkg/log"
	"time"

	"mobingi/ocean/pkg/services/tencent"

	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
)

type Node struct {
	Client *kubernetes.Clientset
}

func (n *Node) GetNode() (*v1.NodeList, error) {
	nodes, err := n.Client.CoreV1().Nodes().List(metav1.ListOptions{})
	if err != nil {
		return nil, err
	}
	return nodes, nil
}

func (n *Node) DeleteNode(name string) error {
	err := n.Client.CoreV1().Nodes().Delete(name, &metav1.DeleteOptions{})
	if err != nil {
		return err
	}
	return nil
}

func (n *Node) GetUnhealthyNodeNum() (int64, error) {
	var num int64 = 0
	nodes, err := n.GetNode()
	if err != nil {
		return num, err
	}
	for _, node := range nodes.Items {
		for _, condition := range node.Status.Conditions {
			if condition.Type == "Ready" && condition.Status != "True" {
				num++
			}
		}
	}
	return num, nil
}

func (n *Node) NewUnhealthyNodeTimer() {
	tencent.Init()
	client := &tencent.InstanceTencent{}
	var lastNum int64 = 0
	lastTime := time.Now()
	timeoutTime := time.Minute * 2

	ticker := time.NewTicker(time.Second * 5)
	go func() {
		for _ = range ticker.C {
			num, err := n.GetUnhealthyNodeNum()
			log.Info(num)
			if err != nil {
				log.Error(err)
			}
			if num > 0 {
				// 小于两分钟
				if lastTime.After(time.Now().Add(-timeoutTime)) {
					if num > lastNum {
						// TODO 将实例id 对应集群
						res, err := client.CreateInstance(num - lastNum)
						if err != nil {
							log.Error(err)
						} else {
							log.Info(res)
							lastNum = num
							lastTime = time.Now()
						}
					}
				} else {
					res, err := client.CreateInstance(num)
					if err != nil {
						log.Error(err)
					} else {
						log.Info(res)
						lastNum = num
						lastTime = time.Now()
					}
				}
			}
		}
	}()

}
