package certs

import (
	"errors"
	"fmt"
	"mobingi/ocean/pkg/config"
	"net"

	certutil "k8s.io/client-go/util/cert"
)

func makeAltNamesMutator(f func(*config.Config) (*certutil.AltNames, error)) configMutatorsFunc {
	return func(c *config.Config, cc *certutil.Config) error {
		altNames, err := f(c)
		if err != nil {
			return err
		}
		cc.AltNames = *altNames
		return nil
	}
}

func getAPIServerAltNames(cfg *config.Config) (*certutil.AltNames, error) {
	advertiseAddress := net.ParseIP(cfg.AdvertiseAddress)
	if advertiseAddress == nil {
		return nil, errors.New("must have advertiseAddress")
	}

	// TODO fix up
	/*
		_, svcSubnet, err := net.ParseCIDR(cfg.Networking.ServiceSubnet)
		if err != nil {
			return nil, fmt.Errorf("error parsing CIDR %q", cfg.Networking.ServiceSubnet)
		}*/

	/*internalAPIServerVirtualIP, err := ipallocator.GetIndexdIP(svcSubnet, 1)
	if err != nil {
		return nil, fmt.Errorf("unable to get first IP address from the given CIDR (%s), error:%s", cfg.Networking.ServiceSubnet, err.Error())
	}*/

	altNames := &certutil.AltNames{
		DNSNames: []string{
			//cfg.NodeRegistration.Name,
			"kubernetes",
			"kubernetes.default",
			"kubernetes.default.svc",
			fmt.Sprintf("kubernetes.default.svc.%s", cfg.Networking.DNSDomain),
		},
		IPs: []net.IP{
			//internalAPIServerVirtualIP,
			advertiseAddress,
		},
	}

	//TODO fix up

	return altNames, nil
}

func getEtcdAltNames(cfg *config.Config) (*certutil.AltNames, error) {
	advertiseAddress := net.ParseIP(cfg.AdvertiseAddress)
	if advertiseAddress == nil {
		return nil, errors.New("must have advertiseAddress")
	}

	altNames := &certutil.AltNames{
		IPs: []net.IP{
			advertiseAddress,
			net.ParseIP("0.0.0.0"),
		},
	}

	//TODO fix up

	return altNames, nil
}
