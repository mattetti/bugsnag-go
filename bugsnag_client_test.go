package bugsnag

import (
	"errors"
	"os"
	"testing"
)

func TestSendToLiveAPI(t *testing.T) {
	Verbose = true
	APIKey = os.Getenv("BUGSNAG_APIKEY")
	AppVersion = "1.2.3"
	OSVersion = "3.2.1"
	if err := Notify(errors.New("This is a test"), "12345"); err != nil {
		t.Fatal(err)
	}
}
