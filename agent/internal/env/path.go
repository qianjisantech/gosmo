package env

import (
	"fmt"
)

func GetGorPath() (string, error) {
	switch GetOS() {
	case "windows":
		return "./cmd/gor.exe", nil
	case "linux":
		return "./cmd/gor", nil
	//case "mac":
	//	return "./cmd/linux/gor", nil
	default:
		return "", fmt.Errorf("unsupported OS")
	}
}
