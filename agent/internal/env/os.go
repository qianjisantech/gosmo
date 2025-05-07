package env

import "runtime"

func GetOS() string {
	switch runtime.GOOS {
	case "windows":
		return "windows"
	case "linux":
		return "linux"
	case "darwin":
		return "mac"
	default:
		return "unknown"
	}
}
