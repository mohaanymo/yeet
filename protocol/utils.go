package protocol

import (
	"os"
)

func IsPathExist(filePath string) (bool, error) {
	_, err := os.Stat(filePath)
	if err!=nil {
		if err == os.ErrNotExist {
			return false, nil
		}

		return false, err
	}

	return true, nil

}
