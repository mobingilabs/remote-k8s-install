package dependence

import (
	"mobingi/ocean/pkg/tools/machine"
)

// get from fileServer
func getMasterBinCommands() machine.CommandList {
	cl := machine.CommandList{}

	// TODO from config
	curlCmd := "curl -L localhost:3232/master"
	curlCheck := func(output string) bool {
		return true
	}
	cl.Add(curlCmd, curlCheck)

	tarCmd := "tar -zcvf xxx -o dasdas"
	tarCheck := func(output string) bool {
		return true
	}
	cl.Add(tarCmd, tarCheck)

	return cl
}

func getNodeBinCommands() machine.CommandList {
	cl := machine.CommandList{}

	// TODO from config
	curlCmd := "curl -L localhost:3232/master"
	curlCheck := func(output string) bool {
		return true
	}
	cl.Add(curlCmd, curlCheck)

	tarCmd := "tar -zcvf xxx -o dasdas"
	tarCheck := func(output string) bool {
		return true
	}
	cl.Add(tarCmd, tarCheck)

	return cl
}

func getCniBinCommands() machine.CommandList {
	cl := machine.CommandList{}

	// TODO from config
	curlCmd := "curl -L localhost:3232/master"
	curlCheck := func(output string) bool {
		return true
	}
	cl.Add(curlCmd, curlCheck)

	tarCmd := "tar -zcvf xxx -o dasdas"
	tarCheck := func(output string) bool {
		return true
	}
	cl.Add(tarCmd, tarCheck)

	return cl
}