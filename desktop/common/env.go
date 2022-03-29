package common

import (
	"os"
	"path/filepath"
)

var (
	appPath = ""
)

func GetAppPath() string {
	if appPath == "" {
		appPath = filepath.Dir(os.Args[0])
	}
	return appPath
}

func GetBinPath() string {
	return filepath.Join(GetAppPath(), "bin")
}
