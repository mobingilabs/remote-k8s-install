package service

import (
	"fmt"
	"strings"
)

const (
	EtcdClientPort = 2379
)

func GetEtcdServers(ips []string) string {
	urls := make([]string, 0, len(ips))
	for _, v := range ips {
		url := fmt.Sprintf("https://%s:%d", v, EtcdClientPort)
		urls = append(urls, url)
	}

	return strings.Join(urls, ",")
}
