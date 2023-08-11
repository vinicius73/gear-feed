package support

import (
	"os"
	"path/filepath"
)

func GetBinDirPath() string {
	execPath, _ := os.Executable()

	return filepath.Dir(execPath)
}

func DirMustExist(dirPath string) error {
	if _, err := os.Stat(dirPath); os.IsNotExist(err) {
		return os.MkdirAll(dirPath, os.ModePerm)
	}

	return nil
}
