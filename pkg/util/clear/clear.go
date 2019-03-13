package clear

import (
	"os"
)

const workDir = "/etc/kubernetes"

func Clean() {
	os.RemoveAll(workDir)
}
