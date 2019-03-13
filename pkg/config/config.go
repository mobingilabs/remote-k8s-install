package config

import ()

const (
	workDir = "/etc/kubernetes"
	pkiDir  = "/etc/kubernetes/pki"
)

type MachineRole int

const (
	MasterRole MachineRole = iota
	NodeRole
)

type Config struct {
	PKIDir           string
	WorkDir          string
	AdvertiseAddress string
	ClusterName      string
	Machines         []Machine

	Networking Networking
}

// TODO more ssh auth support
type Machine struct {
	Role     MachineRole
	Addr     string
	User     string
	Password string
}

func NewConfig() *Config {
	return &Config{
		PKIDir:           pkiDir,
		AdvertiseAddress: "192.168.1.74",
		WorkDir:          workDir,
		ClusterName:      "kubernetes",
		Machines: []Machine{
			Machine{
				Role:     MasterRole,
				Addr:     "192.168.1.74",
				User:     "root",
				Password: "947337",
			},
		},
	}
}

// TODO more than one machine get
func (c *Config) GetMasterMachine() *Machine {
	for _, v := range c.Machines {
		if v.Role == MasterRole {
			return &v
		}
	}

	return nil
}

type Networking struct {
	ServiceSubnet string
	PodSubnet     string
	DNSDomain     string
}
