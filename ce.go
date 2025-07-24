//go:build !pro

package gosmo

import (
	"fmt"
)

// PRO this value indicates if goreplay is running in PRO mode.
var PRO = false

func SettingsHook(settings *AppSettings) {
	if settings.RecognizeTCPSessions {
		settings.RecognizeTCPSessions = false
		fmt.Println("[ERROR] TCP session recognition is not supported in the open-source version of GoReplay")
	}
}
