package support

import (
	"os"
	"path/filepath"
)

func GetBinDirPath() string {
	execPath, _ := os.Executable()

	return filepath.Dir(execPath)
}
