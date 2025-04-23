//go:build recorder
// +build recorder

package sharingan

import (
	"log"
	"os"

	"github.com/qianjisantech/gosmo/plugins"
	"github.com/qianjisantech/gosmo/recorder"
	"github.com/qianjisantech/gosmo/recorder/koala/hook"
	"github.com/qianjisantech/gosmo/recorder/koala/logger"
	"github.com/qianjisantech/gosmo/recorder/koala/sut"
)

// GetCurrentGoRoutineID get current goroutineID incase SetDelegatedFromGoRoutineID
func GetCurrentGoRoutineID() int64 {
	return recorder.GetCurrentGoRoutineID()
}

// SetDelegatedFromGoRoutineID should be used when this goroutine is doing work for another goroutine
func SetDelegatedFromGoRoutineID(gID int64) {
	recorder.SetDelegatedFromGoRoutineID(gID)
}

func init() {
	if os.Getenv("RECORDER_ENABLED") != "true" {
		return
	}

	// init logger
	logger.Init()

	// init plugin && start recorder
	plugins.InitRecorderPlugin()
	plugins.StartRecorder()

	// start hook
	hook.Start()

	// start gc
	sut.StartGC()

	// log
	log.Println("mode", "=====recorder=====")
}
