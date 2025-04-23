//go:build !recorder && !recorder_grpc && !replayer
// +build !recorder,!recorder_grpc,!replayer

package gosmo

// GetCurrentGoRoutineID get current goroutineID incase SetDelegatedFromGoRoutineID
func GetCurrentGoRoutineID() int64 {
	return 0
}

// SetDelegatedFromGoRoutineID should be used when this goroutine is doing work for another goroutine
func SetDelegatedFromGoRoutineID(gID int64) {
}

func init() {
}
