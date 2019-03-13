package dir

import (
	"os"
)

func MkdirAllIfNotExists(dir string) error {
	_, err := os.Stat(dir)
	if err != nil {
		return nil
	}
	if os.IsNotExist(err) {
		if err := os.MkdirAll(dir, 0777); err != nil {
			return err
		}
	}

	return nil
}
