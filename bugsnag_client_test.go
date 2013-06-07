package bugsnag

import (
	"errors"
	"os"
	"testing"
)

func TestNotify(t *testing.T) {
	// Configure bugsnap
	Verbose = true
	APIKey = os.Getenv("BUGSNAG_APIKEY")
	AppVersion = "1.2.3"
	OSVersion = "3.2.1"

	// Notify about an error
	e := errors.New("This is a test")
	if err := Notify(e); err != nil {
		t.Fatal(err)
	}
}

func TestSetMetaDataBeforeNotify(t *testing.T) {
	// Configure bugsnap
	APIKey = os.Getenv("BUGSNAG_APIKEY")
	Verbose = true

	// Notify about another error, with more details
	e := errors.New("This is another test")
	values := map[string]interface{}{
		"account_id": 5555,
		"user_agent": "ie4",
	}
	if err := New(e).SetUserID("12345").SetMetaData("user_info", values).Notify(); err != nil {
		t.Fatal(err)
	}
}

func TestAddMetaDataBeforeNotify(t *testing.T) {
	// Configure bugsnap
	APIKey = os.Getenv("BUGSNAG_APIKEY")
	Verbose = true

	// Notify about another error, with more details
	e := errors.New("Another error")
	if err := New(e).SetUserID("12345").AddMetaData("user_info", "name", "mr. Fu").Notify(); err != nil {
		t.Fatal(err)
	}
}

func TestNewNotify(t *testing.T) {
	// Configure bugsnap
	APIKey = os.Getenv("BUGSNAG_APIKEY")
	Verbose = true

	// Notify about another error, with more details
	e := errors.New("One more error")
	if err := New(e).Notify(); err != nil {
		t.Fatal(err)
	}
}
