package util

import "runtime"

func GetDroneManifestExecCmd() string {
	if runtime.GOOS == "windows" {
		return "C:/bin/drone-manifest.exe"
	}

	return "drone-manifest"
}
