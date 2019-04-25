package certs

import "fmt"

func PathForCert(baseName string) string {
	return fmt.Sprintf("%s.crt", baseName)
}

func PathForKey(baseName string) string {
	return fmt.Sprintf("%s.key", baseName)
}

func PathForPub(baseName string) string {
	return fmt.Sprintf("%s.pub", baseName)
}
