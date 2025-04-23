package gosmo_test

import (
	"net/http"
)

func Example() {
	doneChan := make(chan bool)

	go func(delegatedID int64) {
		gosmo.SetDelegatedFromGoRoutineID(delegatedID)
		defer gosmo.SetDelegatedFromGoRoutineID(0)

		http.Get("http://127.0.0.1:8888")

		doneChan <- true
	}(gosmo.GetCurrentGoRoutineID())

	<-doneChan
}
