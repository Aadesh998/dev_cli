package utils

import (
	"log"
	"os"
	"runtime"
)

func GetDirConfigPath() string {
	configPath, err := os.UserConfigDir()
	if err != nil {
		log.Printf("Failed to get config path on %s: %v", runtime.GOOS, err)
		return ""
	}
	return configPath
}
